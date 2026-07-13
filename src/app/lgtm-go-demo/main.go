package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	// 1. Define the Resource (Identifies your service)
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("lgtm-demo-service-go"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// 2. Setup Tracing
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(), // Insecure is fine for local testing
	)
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	defer tp.Shutdown(ctx)
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("demo-tracer")

	// 3. Setup Metrics
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create metric exporter: %v", err)
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	defer meterProvider.Shutdown(ctx)
	otel.SetMeterProvider(meterProvider)
	meter := otel.Meter("demo-meter")

	workCounter, _ := meter.Int64Counter(
		"demo.work.completed",
		metric.WithDescription("Counts work completed"),
	)

	log.Println("Sending data to LGTM stack. Press Ctrl+C to stop.")

	// 4. Generate Telemetry loop
	for {
		// Create a trace (span)
		_, span := tracer.Start(ctx, "do_work", trace.WithAttributes(
			attribute.String("work.type", "demo"),
		))

		// Increment the metric
		workCounter.Add(ctx, 1)

		// Simulate work
		time.Sleep(1 * time.Second)

		// Add an event to the span and end it
		span.AddEvent("Work completed successfully")
		span.End()

		log.Println("Generated 1 trace and 1 metric increment...")
		time.Sleep(2 * time.Second)
	}
}
