[![Tests](https://github.com/fbuedding/fiware-iot-agent-sdk/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/fbuedding/fiware-iot-agent-sdk/actions/workflows/tests.yml)
**iotagentsdk: Go SDK for IoT Agents**

The `iotagentsdk` package is a Go software development kit (SDK) designed to facilitate interactions with IoT Agents in FIWARE-based systems. FIWARE is an open-source initiative aimed at providing a framework for the development of smart solutions in various domains, including smart cities, industrial IoT, and agriculture.

**Features:**

- **Config Group Management:** Allows you to manage configuration groups for devices, including creating, reading, updating, and deleting (CRUD) operations.
  
- **Device Management:** Provides functionalities for managing devices, such as reading device information, checking device existence, listing devices, creating devices, updating device information, and deleting devices.

**How to Use:**

1. **Import the Package:** Import the `iotagentsdk` package into your Go project:

   ```go
   import "github.com/fbuedding/fiware-iot-agent-sdk"
   ```

2. **Initialize the IoTA Client:** Create an instance of the `IoTA` struct, providing the necessary configuration parameters such as host, port, and timeout.

   ```go
   iotAgent := iotagentsdk.NewIoTA("localhost", 4041, 5000)
   ```

3. **Interact with IoT Agents:** Use the methods provided by the `iotagentsdk` package to perform operations such as managing configuration groups and devices.

   ```go
   // Example: Read device information
   deviceID := "my-device"
   deviceInfo, err := iotAgent.ReadDevice(myFiwareService, iotagentsdk.DeciveId(deviceID))
   if err != nil {
       log.Fatal("Error reading device:", err)
   }
   fmt.Println("Device information:", deviceInfo)
   ```

**Compatibility:**

- The `iotagentsdk` package is compatible with Go versions 1.11 and higher.

**Contributing:**

Contributions to the `iotagentsdk` package are welcome! Feel free to submit bug reports, feature requests, or pull requests on GitHub.

**License:**

This package is licensed under the MIT License. See the LICENSE file for details.

For more information and detailed usage instructions, please refer to the [documentation](https://github.com/yourusername/iotagentsdk).

**Acknowledgments:**

This package is developed based on the FIWARE architecture and specifications. I acknowledge the contributions of the FIWARE community and developers.
