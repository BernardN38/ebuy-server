package com.ebuy.mediaservice.messaging;

import java.io.Serializable;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

import lombok.Data;

@Data
public class ImageUploadedMessage implements Serializable {
    private Long mediaId;
    private UUID externalIdFull;
    private UUID externalIdCompressed;
    private String contentType;
    private Long userId;

    @JsonCreator
    public ImageUploadedMessage(
            @JsonProperty("mediaId") Long mediaId,
            @JsonProperty("externalIdFull") UUID externalIdFull,
            @JsonProperty("externalIdCompressed") UUID externalIdCompressed,
            @JsonProperty("contentType") String contentType,
            @JsonProperty("userId") Long userId) {
        this.mediaId = mediaId;
        this.externalIdFull = externalIdFull;
        this.externalIdCompressed = externalIdCompressed;
        this.contentType = contentType;
        this.userId = userId;
    }
}
