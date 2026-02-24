# Changelog

All notable changes to this project will be documented in this file.

## [v1.2.0] - 2026-02-24

### Added
- ✨ Added OSS (Object Storage Service) support for SSL certificate deployment to custom domains
- ✨ Added OSS API integration with certificate update functionality for CNAME records
- ✨ Added conditional service deployment - services are now only deployed if configured
- ✨ Added support for certificate deployment using certificate content when certificate ID is not available

### Changed
- ♻️ Refactored deploy logic to support conditional service deployment based on configuration
- ♻️ Modified CAS integration to use certificate content directly instead of file paths
- ♻️ Updated CDN integration to handle both CAS certificate IDs and certificate content
- ♻️ Enhanced SLB integration to support both certificate ID and certificate content upload
- ♻️ Improved error logging and success messaging throughout deployment process

### Fixed
- 🐛 Fixed missing nil checks that could cause panic if service configurations are not present
- 🐛 Fixed redundant return statements in deploy function that prevented complete deployment
- 🐛 Fixed typo in SLB logging ("htttps" corrected to "https")
- 🐛 Fixed typo in CDN logging ("certicating" corrected to "certifying")

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
Unique failure occurred: oldString not found in content

## [v1.0.0] - 2026-02-23

### Added
- Initial release with CAS (Certificate Authority Service) support
- Basic certificate deployment functionality for Alibaba Cloud
- Configuration system for multiple accounts
- Integration with acme.sh for automated certificate management