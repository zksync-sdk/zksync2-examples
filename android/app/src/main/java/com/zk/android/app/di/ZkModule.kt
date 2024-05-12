package com.zk.android.app.di

import com.zk.android.app.BuildConfig
import com.zk.android.app.network.ZkConfig
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.zksync.protocol.ZkSync
import org.web3j.crypto.Credentials
import org.web3j.protocol.Web3j
import org.web3j.protocol.http.HttpService
import javax.inject.Singleton


@Module
@InstallIn(SingletonComponent::class)
class ZkModule {

    @Provides
    @Singleton
    fun provideConfig(): ZkConfig {
        return ZkConfig(
            BuildConfig.RPC_ETH_URL,
            BuildConfig.ZK_ERA_URL,
            BuildConfig.ZK_ERA_SOCKET,
            BuildConfig.PRIVATE_KEY
        )
    }

    @Provides
    @Singleton
    fun provideWeb3j(config: ZkConfig): Web3j {
        return Web3j.build(HttpService(config.rpc))
    }

    @Provides
    @Singleton
    fun provideZkSync(config: ZkConfig): ZkSync {
        return ZkSync.build(HttpService(config.url))
    }

    @Provides
    @Singleton
    fun provideCredentials(config: ZkConfig): Credentials {
        return Credentials.create(config.privateKey)
    }
}