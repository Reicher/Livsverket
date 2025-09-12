# Livsverket

Livsverket is a local-only Android app that lets you explore the tree of life and build a personal collection of organisms. The app works entirely offline, storing a simplified biological taxonomy in a local JSON or SQLite database. Users can navigate through domains, kingdoms, species and more, choose to collect any taxonomic set, and report sightings.

## Features
- Navigate the biological taxonomy starting from the three domains: **Bacteria**, **Archaea**, and **Eukaryotes**.
- Drill down or move back up the tree by tapping on sets.
- Collect any set or mark it as seen.
- View your collection with stats on how many child sets you have collected.

## Getting Started

### Prerequisites
- [Android Studio](https://developer.android.com/studio) (Arctic Fox or newer) on Ubuntu
- Android SDK 26+
- Java 11+ (bundled with Android Studio)

### Clone the Repository
```bash
git clone https://github.com/your-username/Livsverket.git
cd Livsverket
```

### Open in Android Studio
1. Launch Android Studio.
2. Select **Open an Existing Project**.
3. Choose the `Livsverket` directory.
4. Allow Android Studio to sync and download required dependencies.

### Build and Run
1. Connect an Android device with USB debugging enabled, or start an emulator from the **AVD Manager**.
2. In Android Studio, click the **Run** button (▶) or use `Shift + F10`.
3. The app will compile and install on the selected device.

## Enabling USB Debugging on Android
1. Open the device **Settings**.
2. Navigate to **About phone** and tap **Build number** seven times to enable Developer Options.
3. Return to **Settings** → **System** → **Developer options**.
4. Enable **USB debugging**.
5. Connect the device via USB and authorize the computer when prompted.

## Troubleshooting
- **Gradle sync fails**: Ensure you have an active internet connection on first build and that the required SDK components are installed via the **SDK Manager**.
- **Device not detected**: Verify that USB debugging is enabled and that `adb devices` lists your device.
- **Build errors**: Clean the project with **Build → Clean Project** and try again. Check that all required modules and plugins are installed.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
