# Object storage - primary integration owner: Waheed

Implement an S3-compatible adapter and narrow interface when media work starts. No bucket or provider credentials are configured in the scaffold.

Use restricted signed uploads, ownership checks, allowed image types, size limits, image decoding/validation, random object keys, and explicit completion. Keep pending uploads/private exports private. Publish only validated product assets through the CDN. Do not expose object storage credentials to browsers or logs.
