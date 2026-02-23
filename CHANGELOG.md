# Changelog

All notable changes to this project will be documented in this file.

## [v1.1.0] - 2026-02-23

### Added
- ✨ Added CDN (Content Delivery Network) support for SSL certificate deployment
- ✨ Added SLB (Server Load Balancer) support for SSL certificate deployment
- ✨ Created shared utility functions for Alibaba Cloud API client initialization
- ✨ Added domain validation to prevent injection attacks and malformed inputs
- ✨ Implemented pagination bounds checking to prevent excessive API calls
- ✨ Added comprehensive documentation for all supported services

### Changed
- 🐛 Fixed inconsistent error messages that incorrectly referenced CDN when performing SLB operations
- ♻️ Refactored certificate loading to use centralized model with validation
- 🔧 Updated Makefile to be more security-conscious when installing binaries
- 📝 Improved README.md with proper security warnings and clearer instructions

### Fixed
- 🐛 Fixed typo in error messages to ensure accurate debugging information
- 🐛 Implemented proper error handling for Alibaba Cloud client initialization to prevent potential panics
- 🐛 Fixed redundant export statements in deployment script example
- 🐛 Corrected RegionId parameter usage in CAS integration

## [v1.0.0] - 2026-02-23

### Added
- Initial release with CAS (Certificate Authority Service) support
- Basic certificate deployment functionality for Alibaba Cloud
- Configuration system for multiple accounts
- Integration with acme.sh for automated certificate management