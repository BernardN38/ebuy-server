package com.ebuy.mediaservice.service;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import org.springframework.web.multipart.MultipartFile;

import com.ebuy.mediaservice.entities.UserMedia.UserMedia;

import io.minio.GetObjectResponse;

public interface MediaService {

    List<UserMedia> GetAll() throws Exception;

    Optional<GetObjectResponse> GetImage(String imageId) throws Exception;

    UUID CreateUserMedia(MultipartFile image, Long userId) throws Exception;
}