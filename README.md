# QeVR

<p align="center">
  <img src="logo.jpeg" alt="QeVR Logo" width="200">
</p>



QeVR is a tool designed to help upload vulnerability scans to Tipping Point SMS server for profile tuning. It streamlines the process of transferring vulnerability data from Qualys/Rapid7 to SMS, enhancing the efficiency of security management workflows.

## Features

- Import security scan data in CSV format with heuristics detection of columns with IP addresses and CVE IDs
- Filter vulnerabilities data by CIDR networks
- Export data in a format compatible with Tipping Point SMS
- Upload to Tipping Point SMS directly
- User-friendly graphical interface for easy operation

## Installation

Download the latest release from [GitHub Releases](https://github.com/mpkondrashin/qevr/releases/latest) and extract to any appropriate folder.

## Usage

1. Launch the QeVR application
2. Follow the wizard-style interface to:
   - Select the source of vulnerability scan data
   - Configure output settings
   - Process and export the data

## Requirements

QeVR supports Windows (x64 platform) and macOS (ARM/x64 platforms)

## Configuration

QeVR saves all parameters in config.yaml file aside to QeVR executable upon successful run.

## License

QeVR is distributed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contributing

Contributions to QeVR are welcome! Please feel free to submit pull requests, report bugs, or suggest features.

## Support

For support, please open an issue on the [GitHub repository](https://github.com/mpkondrashin/qevr) or contact the maintainer directly.

## Acknowledgements

QeVR is developed and maintained by Mikhail Kondrashin (mkondrashin@gmail.com).
