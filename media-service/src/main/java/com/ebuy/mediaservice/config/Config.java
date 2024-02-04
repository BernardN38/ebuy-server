package com.ebuy.mediaservice.config;

import java.io.IOException;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.ebuy.mediaservice.fileManager.FileManager;

@Configuration
public class Config {

    @Bean
    public FileManager fileUploader() throws InvalidKeyException,
            NoSuchAlgorithmException, IOException {
        return new FileManager("media-service");
    }

}
