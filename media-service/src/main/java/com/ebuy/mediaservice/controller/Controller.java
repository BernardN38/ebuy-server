package com.ebuy.mediaservice.controller;

import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.util.List;
import java.util.UUID;

import org.checkerframework.checker.units.qual.m;
import org.hibernate.mapping.Map;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.InputStreamResource;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.security.core.Authentication;

import com.ebuy.mediaservice.dto.MediaUploadResponse;
import com.ebuy.mediaservice.entities.UserMedia.UserMedia;
import com.ebuy.mediaservice.jwtFilter.User;
import com.ebuy.mediaservice.service.MediaService;
import com.fasterxml.jackson.core.util.ByteArrayBuilder;

import io.minio.GetObjectResponse;

@RestController
public class Controller {

    private final MediaService mediaService;

    public Controller(MediaService mediaService) {
        this.mediaService = mediaService;
    }

    @GetMapping("/api/v1/media/health")
    String health() {
        return "media-service up and running";
    }

    @PostMapping("api/v1/media")
    public ResponseEntity<MediaUploadResponse> handleImageUpload(
            @RequestParam("image") MultipartFile image) throws Exception {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        User user = (User) authentication.getPrincipal();

        UUID mediaId = mediaService.CreateUserMedia(image, Long.valueOf(user.getUserId()));
        return ResponseEntity
                .created(URI.create("/api/v1/media/" + mediaId)).body(new MediaUploadResponse(mediaId));
    }

    @GetMapping("api/v1/media/{id}")
    public ResponseEntity<ByteArrayResource> GetImage(@PathVariable String id) throws IOException {
        ByteArrayResource media;
        try {
            GetObjectResponse objResp = mediaService.GetImage(id).orElseThrow();
            media = new ByteArrayResource(objResp.readAllBytes());
        } catch (Exception e) {
            return ResponseEntity.badRequest().build();
        }
        return ResponseEntity.ok()
                .contentType(MediaType.IMAGE_JPEG)
                .body(media);
    }

    @GetMapping("api/v1/media")
    public ResponseEntity<List<UserMedia>> GetAll() throws Exception {
        List<UserMedia> allUserMedia = mediaService.GetAll();
        return ResponseEntity.ok().body(allUserMedia);
    }
}
