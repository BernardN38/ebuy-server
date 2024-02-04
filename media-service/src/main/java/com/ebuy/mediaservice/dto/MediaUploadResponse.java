package com.ebuy.mediaservice.dto;

import java.util.UUID;

import lombok.Data;

@Data
public class MediaUploadResponse {
    private Long mediaId;

    public MediaUploadResponse(Long mediaId) {
        this.mediaId = mediaId;
    }
}
