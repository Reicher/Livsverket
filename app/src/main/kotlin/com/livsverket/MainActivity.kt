package com.livsverket

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import com.livsverket.ui.theme.LivsverketTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            LivsverketTheme {
                LivsverketApp()
            }
        }
    }
}
