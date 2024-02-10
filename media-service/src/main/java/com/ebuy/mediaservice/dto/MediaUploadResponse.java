package com.ebuy.mediaservice.dto;

import java.util.UUID;

import lombok.Data;

@Data
public class MediaUploadResponse {
    private UUID mediaId;

    public MediaUploadResponse(UUID mediaId) {
        this.mediaId = mediaId;
    }
}
