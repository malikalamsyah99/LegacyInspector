# LegacyInspector

LegacyInspector is a CLI tool for inspecting legacy Apple serial numbers.

The project focuses on older Apple devices that use the legacy 11-character and 12-character serial number format.

## Features

- Detect 11-character legacy serial numbers
- Detect 12-character legacy serial numbers
- Detect randomized serial numbers
- Decode production year
- Decode production week
- Decode production and product codes
- Identify Apple devices from product codes
- Identify model identifier and model year when available
- Resolve ambiguous year codes using known device information
- Identify manufacturing location when available

## Usage

Run LegacyInspector with a serial number:

    go run ./cmd/legacyinspector SERIAL

Example:

    go run ./cmd/legacyinspector C3WF2E79DCP9

Example output:

    LegacyInspector
    ================================
    Serial         : C3WF2E79DCP9
    Format         : Legacy 12-character
    Valid          : true
    Location Code  : C3W
    Location       : Unknown
    Year Code      : F
    Production Year : 2011
    Half           : 1
    Week Code      : 2
    Week           : 2
    Production     : E79
    Product Code   : DCP9
    Device         : iPod touch (4th generation)
    Reason         : Valid Apple legacy 12-character structure. Product recognized.
    ================================

## Year Resolution

Legacy serial year codes can represent different decades.

For example, `F` can represent different production years depending on the device.

LegacyInspector uses the identified device/model information to resolve the year instead of always assuming a single year.

Example:

    C02FN15CXXXX

    Device          : MacBook Pro (13-inch, 2020, Apple M1)
    Model Identifier: MacBookPro17,1
    Model Year      : 2020
    Production Year : 2021

While:

    C3WF2E79DCP9

    Device         : iPod touch (4th generation)
    Production Year: 2011

## Data

LegacyInspector currently uses several datasets:

- `products.json` — Apple product codes
- `macserial_products.json` — additional product code data
- `models.json` — model data generated from OpenCore
- `model_overrides.json` — manually verified model mappings
- `overrides.json` — product overrides
- `locations.json` — legacy manufacturing location codes

The OpenCore source used to generate `models.json` is not included in the repository.

## Project Structure

    LegacyInspector/
    ├── cmd/
    │   └── legacyinspector/
    │       └── main.go
    ├── internal/
    │   └── decoder/
    │       ├── data/
    │       │   ├── locations.json
    │       │   ├── macserial_products.json
    │       │   ├── model_overrides.json
    │       │   ├── models.json
    │       │   ├── overrides.json
    │       │   └── products.json
    │       ├── decoder.go
    │       ├── location.go
    │       ├── location_loader.go
    │       ├── model_loader.go
    │       ├── product.go
    │       ├── product_loader.go
    │       └── validator.go
    └── tools/
        ├── import_locations/
        ├── import_macserials/
        ├── import_models/
        ├── import_products/
        ├── modelinfo.h
        └── modelinfo_autogen.h

## Development

Format and check the project:

    go fmt ./...
    go vet ./...
    go build ./cmd/legacyinspector

Regenerate the model dataset:

    go run ./tools/import_models

The OpenCore source must be available locally under:

    tools/import_models/OpenCorePkg/

This directory is intentionally excluded from Git.

## Current Scope

LegacyInspector currently focuses on legacy Apple serial numbers.

Randomized serial numbers are detected but are not decoded.

Some valid serial numbers may still return `Device: Unknown` when the corresponding product or model information is not available in the current datasets.

## Status

Early development.

The current goal is to build a reliable CLI for identifying and decoding legacy Apple devices from their serial numbers.

More device models, product codes, and manufacturing locations will be added as the dataset is expanded.
