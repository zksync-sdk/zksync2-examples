plugins {
    alias(libs.plugins.androidApplication)
    alias(libs.plugins.jetbrainsKotlinAndroid)
    alias(libs.plugins.daggerHilt)
    alias(libs.plugins.kotlinKapt)
}

android {
    namespace = "com.zk.android.app"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.zk.android.app"
        minSdk = 24
        targetSdk = 33
        versionCode = 1
        versionName = "1.0"
        //TODO Insert your key here
        //buildConfigField("String", "PRIVATE_KEY", "\"<WALLET_PRIVATE_KEY>\"")
    }
    flavorDimensions.add("default")

    productFlavors {
        create("mainnet") {
            applicationIdSuffix = ".debug"
            buildConfigField("String", "RPC_ETH_URL", "\"https://rpc.ankr.com/eth\"")
            buildConfigField("String", "ZK_ERA_URL", "\"https://mainnet.era.zksync.io\"")
            buildConfigField("String", "ZK_ERA_SOCKET", "\"wss://mainnet.era.zksync.io/ws\"")
        }
        create("sepolia") {
            buildConfigField("String", "RPC_ETH_URL", "\"https://rpc.ankr.com/eth_sepolia\"")
            buildConfigField("String", "ZK_ERA_URL", "\"https://sepolia.era.zksync.dev\"")
            buildConfigField("String", "ZK_ERA_SOCKET", "\"wss://sepolia.era.zksync.dev/ws\"")
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(getDefaultProguardFile("proguard-android-optimize.txt"), "proguard-rules.pro")
        }
    }
    kotlinOptions {
        jvmTarget = "1.8"
    }
    buildFeatures {
        viewBinding = true
        buildConfig = true
    }
}

dependencies {

    //Zk
    implementation(libs.zksync)
    implementation(libs.web3j.core)
    implementation(libs.web3j.crypto)

    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.core.splashscreen)
    implementation(libs.androidx.appcompat)
    implementation(libs.androidx.constraintlayout)
    implementation(libs.androidx.navigation.fragment.ktx)
    implementation(libs.androidx.navigation.ui.ktx)
    implementation(libs.lifecycle.viewmodel.ktx)
    implementation(libs.lifecycle.runtime.ktx)

    implementation(libs.material)
    implementation(libs.orbit.viewmodel)
    implementation(libs.lottie)

    // Dependency injection
    implementation(libs.hilt.android)
    kapt(libs.hilt.compiler)
    kapt(libs.androidx.hilt.compiler)
    implementation(libs.hilt.navigation.fragment)
}