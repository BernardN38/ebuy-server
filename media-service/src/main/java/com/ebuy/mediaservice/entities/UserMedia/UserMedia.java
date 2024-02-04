package com.ebuy.mediaservice.entities.UserMedia;

import java.time.Instant;
import java.util.Date;
import java.util.UUID;

import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.SourceType;
import org.hibernate.annotations.UpdateTimestamp;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.Temporal;
import jakarta.persistence.TemporalType;
import lombok.Getter;
import lombok.Setter;

@Entity
@Getter
@Setter
public class UserMedia {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    @Column(name = "user_id", nullable = false)
    private Long userId;
    @Column(name = "media_id", unique = true, nullable = false)
    private UUID mediaId;
    @Column(name = "media_id_compressed", unique = true)
    private UUID mediaIdCompressed;
    @CreationTimestamp(source = SourceType.DB)
    private Instant createdOn;
    @UpdateTimestamp(source = SourceType.DB)
    private Instant lastUpdatedOn;
    @Column(name = "compression_status", nullable = false)
    private Boolean compression_status = false;

    protected UserMedia() {
    }

    public UserMedia(Long userId, UUID mediaId, UUID mediaIdCompressed) {
        this.userId = userId;
        this.mediaId = mediaId;
        this.mediaIdCompressed = mediaIdCompressed;
    }

    @Override
    public String toString() {
        return String.format(
                "UserMedia[id=%d, userId='%s', mediaId='%s']",
                id, userId, mediaId);
    }
}
