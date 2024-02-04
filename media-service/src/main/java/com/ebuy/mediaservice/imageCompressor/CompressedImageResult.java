package com.ebuy.mediaservice.imageCompressor;

import java.io.ByteArrayOutputStream;

public class CompressedImageResult {
    private final ByteArrayOutputStream compressedByteArray;
    private final long size;

    public CompressedImageResult(ByteArrayOutputStream byteArr, long size) {
        this.compressedByteArray = byteArr;
        this.size = size;
    }

    public ByteArrayOutputStream getCompressedByteArrayOutputStream() {
        return compressedByteArray;
    }

    public long getSize() {
        return size;
    }
}
