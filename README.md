# Hyper-V Detection Tool

A simple Go program that detects whether the current Windows system is running inside a Hyper-V virtual machine environment using CPUID instruction.

## How It Works

The program uses the CPUID instruction to:
1. Check for the Hyper-V vendor string ("Microsoft Hv")
2. Verify specific feature flags that indicate a Hyper-V environment

## Project Structure

- `main.go` - Main program logic and Hyper-V detection implementation
- `cpuid_amd64.s` - Assembly code for CPUID instruction

## Technical Details

The detection is performed by:
1. Using CPUID with leaf `0x40000000` to check for Hyper-V vendor signature
2. Using CPUID with leaf `0x40000003` to verify specific feature flags

## Output

The program will output one of two messages:
- "Running in Hyper-V environment"
- "Not running in Hyper-V environment"
