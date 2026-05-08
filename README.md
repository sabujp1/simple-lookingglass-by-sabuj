# Looking Glass - ISP Network Diagnostics Platform

[![GitHub Stars](https://img.shields.io/github/stars/sabujp1/simple-lookingglass-by-sabuj)](https://github.com/sabujp1/simple-lookingglass-by-sabuj/stargazers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

> A full-stack web application for performing network diagnostics (ping, traceroute, BGP lookup) across multiple routers from various vendors (MikroTik, Juniper, Cisco, Huawei).

[**🚀 Live Demo**](http://localhost:3000) | [**📦 Installation**](#quick-start) | [**📖 Documentation**](docs/) | [**🐛 Report Issue**](https://github.com/sabujp1/simple-lookingglass-by-sabuj/issues)

## 🔗 Links

- **GitHub Repository**: https://github.com/sabujp1/simple-lookingglass-by-sabuj.git
- **Report a Bug**: https://github.com/sabujp1/simple-lookingglass-by-sabuj/issues
- **Pull Requests**: https://github.com/sabujp1/simple-lookingglass-by-sabuj/pulls

## ✨ Features

- **Multi-Vendor Support**: SSH-based connection to MikroTik, Juniper, Cisco, and Huawei routers
- **Network Diagnostics**: Execute ping, traceroute, and BGP lookup queries through a web interface
- **Real-time Streaming**: Live query results via WebSocket connections
- **Role-Based Access Control**: Admin, Operator, and User roles with appropriate permissions
- **Audit Logging**: Complete audit trail of all user actions
- **Responsive UI**: Modern Next.js frontend with dark/light mode support
- **Caching**: Redis-based caching for improved performance
- **Job Queue**: Background processing for long-running queries

## Architecture

## Quick Start

### 1. Clone the repository
```bash
git clone https://github.com/sabujp1/simple-lookingglass-by-sabuj.git
cd simple-lookingglass-by-sabuj
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

For major changes, please open an issue first to discuss what you would like to change.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Sabuj**
- GitHub: [@sabujp1](https://github.com/sabujp1)

## 🙏 Acknowledgments

- Built with Go, Next.js, PostgreSQL, and Redis
- UI components from [shadcn/ui](https://ui.shadcn.com/)
