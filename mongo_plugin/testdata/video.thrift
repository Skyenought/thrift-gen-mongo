namespace go video

struct Video {
    1: i64 Id (go.tag="bson:\"id,omitempty\"")
    2: binary Data (go.tag="bson:\"data,omitempty\"")
}
(
    mongo.InsertVideo = "InsertVideo(ctx context.Context, video *video.Video) (interface{}, error)"
)
