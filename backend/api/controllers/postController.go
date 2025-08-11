package controllers

import (
	"Server/database"
	"Server/models"
	"Server/servergrpc"
	"context"
	"math"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Post
// @Summary create  a new post
// @Description create new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.CreateOrUpdatePost true "post create  deatils"
// @Success 201 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts [post]
func CraetePost(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateOrUpdatePost
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	// start set data
	var post models.PostModel
	post.Creator = c.Locals("userId").(string)
	post.Likes = make([]string, 0)
	post.CreatedAt = time.Now()
	post.Title = body.Title
	post.Message = body.Message
	post.SelectedFile = body.SelectedFile
	//

	var user models.UserModel
	objId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	err := UserSchema.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//
	post.Name = user.Name
	// set data end
	// craete post
	result, err := PostSchema.InsertOne(ctx, &post)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	} else {
		var createdPost *models.PostModel
		query := bson.M{"_id": result.InsertedID}

		PostSchema.FindOne(ctx, query).Decode(&createdPost)
		return c.Status(fiber.StatusCreated).JSON(createdPost)
	}

}

// Get Post
// @Summary Get  a new post
// @Description Get new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post id"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Router /posts/{id} [get]
func GetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}
	return getPostWithComments(c, postID)
}

// Update Post
// @Summary Update  post
// @Description Update post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Param post body models.CreateOrUpdatePost true "update post  deatils"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id} [patch]
func UpdatePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newData models.CreateOrUpdatePost
	if err := c.BodyParser(&newData); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"deatils": err.Error(),
		})
	}

	// authorization start
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost)

	if authPost.Creator != c.Locals("userId").(string) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "You Are Not authorized to update this post.",
		})
	}

	// set data end
	authPost.Title = newData.Title
	authPost.Message = newData.Message
	authPost.SelectedFile = newData.SelectedFile

	// craete post
	// update := bson.M{"title": newData.Title, "message":newData.Message, "selectedFile": newData.SelectedFile}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": authPost.ID}, bson.M{"$set": authPost})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{"data": err.Error()})
	} else {

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": authPost})
	}

}

// GetAllPosts Post
// @Summary Get All Posts
// @Description GetAllPosts with pagination
// @Tags Posts
// @Accept json
// @Produce json
// @Param page query int false "page number"
// @Param id query string true "user id"
// @Success 200 {object} []models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts [get]
func GetAllPosts(c *fiber.Ctx) error {
	var PostSchema = database.DB.Collection("posts")
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	var posts []bson.M // Use bson.M to capture aggregation results

	userid := c.Query("id")
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// get user following list ids and add our user id to it
	MainUserid, _ := primitive.ObjectIDFromHex(userid)
	userSchema.FindOne(ctx, bson.M{"_id": MainUserid}).Decode(&user)
	user.Following = append(user.Following, userid)

	var LIMIT = 2

	filter := bson.M{"creator": bson.M{"$in": user.Following}}

	// get total num of posts
	total, err := PostSchema.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No Posts",
		})
	}

	postCol := database.DB.Collection("posts")

	pipeline := []bson.M{
		{"$match": bson.M{"creator": bson.M{"$in": user.Following}}},
		{"$sort": bson.M{"_id": -1}},
		{"$skip": int64((page - 1) * LIMIT)},
		{"$limit": int64(LIMIT)},
		{"$lookup": bson.M{
			"from": "comments",
			"let":  bson.M{"postId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$postId", "$$postId"}}}},
				{"$sort": bson.M{"createdAt": -1}},
				{"$lookup": bson.M{
					"from": "users",
					"let":  bson.M{"uid": "$userId"},
					"pipeline": []bson.M{
						{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$_id", "$$uid"}}}},
						{"$project": bson.M{"name": 1, "imageUrl": 1}},
					},
					"as": "user",
				}},
				{"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}},
				{"$project": bson.M{"_id": 1, "value": 1, "createdAt": 1, "userId": 1, "user": 1}},
			},
			"as": "comments",
		}},
		// Add this project stage to ensure all fields are included
		{"$project": bson.M{
			"_id":          1,
			"creator":      1,
			"title":        1,
			"message":      1,
			"name":         1,
			"selectedFile": 1,
			"likes":        1,
			"createdAt":    1,
			"comments":     1,
		}},
	}

	cursor, err := postCol.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to aggregate posts",
		})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &posts); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode posts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":          posts,
		"currentPage":   page,
		"numberOfPages": math.Ceil(float64(total) / float64(LIMIT)),
	})
}

// GetPostsUsersBySearch Post
// @Summary Get Posts users by serach query
// @Description get posts adnd users matching the search query
// @Tags Posts
// @Accept json
// @Produce json
// @Param searchQuery query string true "Search query"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/search [get]
func GetPostsUsersBySearch(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var userSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var users []models.UserModel
	var posts []models.PostModel

	//
	filterPost := bson.M{}
	filterUser := bson.M{}

	//
	findOptionsPost := options.Find()
	findOptionsUser := options.Find()

	if search := c.Query("searchQuery"); search != "" {
		// post
		filterPost = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"description": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
		//
		filterUser = bson.M{
			"$or": []bson.M{
				{
					"name": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
				{
					"email": bson.M{
						"$regex": primitive.Regex{
							Pattern: search,
							Options: "i",
						},
					},
				},
			},
		}
	}
	// end
	cursorPosts, _ := PostSchema.Find(ctx, filterPost, findOptionsPost)
	defer cursorPosts.Close(ctx)

	cursorUsers, _ := userSchema.Find(ctx, filterUser, findOptionsUser)
	defer cursorUsers.Close(ctx)
	//

	for cursorUsers.Next(ctx) {
		var user models.UserModel
		cursorUsers.Decode(&user)
		users = append(users, user)
	}

	for cursorPosts.Next(ctx) {
		var post models.PostModel
		cursorPosts.Decode(&post)
		posts = append(posts, post)
	}

	return c.JSON(fiber.Map{
		"user":  users,
		"posts": posts,
	})
}

// Comment Post
// @Summary comment  post
// @Description comment post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Param post body models.ComnmentPost true "comment value"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id}/commentPost [post]
func CommentPost(c *fiber.Ctx) error {
	postCol := database.DB.Collection("posts")
	commentCol := database.DB.Collection("comments")
	userCol := database.DB.Collection("users")
	notificationCol := database.DB.Collection("notifications")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var req models.CreateComment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	// post id
	postID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post id", "details": err.Error()})
	}

	// ensure post exists
	var post models.PostModel
	if err := postCol.FindOne(ctx, bson.M{"_id": postID}).Decode(&post); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	// current user id (from middleware)
	userHex := c.Locals("userId").(string)
	userObjID, _ := primitive.ObjectIDFromHex(userHex)

	// create comment doc
	comment := models.Comment{
		ID:        primitive.NewObjectID(),
		PostID:    postID,
		UserID:    userObjID,
		Value:     req.Value,
		CreatedAt: time.Now(),
	}

	if _, err := commentCol.InsertOne(ctx, comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert comment", "details": err.Error()})
	}

	// create notification
	var user models.UserModel
	if err := userCol.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user); err == nil {
		notification := models.Notification{
			MainUID:   post.Creator,
			TargetID:  postID.Hex(),
			Deatils:   user.Name + " Commented on your Post",
			User:      models.User{Name: user.Name, Avatart: user.ImageUrl},
			CreatedAt: time.Now(),
		}
		res, _ := notificationCol.InsertOne(ctx, notification)
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			notification.ID = oid
			// servergrpc.SendNotification(notification)
		}
	}

	// return the post with live-populated comments using the same function
	return getPostWithComments(c, postID)
}

// Make sure getPostWithComments returns proper structure
func getPostWithComments(c *fiber.Ctx, postID primitive.ObjectID) error {
	postCol := database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"_id": postID}},
		{"$lookup": bson.M{
			"from": "comments",
			"let":  bson.M{"postId": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$postId", "$$postId"}}}},
				{"$sort": bson.M{"createdAt": -1}},
				{"$lookup": bson.M{
					"from": "users",
					"let":  bson.M{"uid": "$userId"},
					"pipeline": []bson.M{
						{"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$_id", "$$uid"}}}},
						{"$project": bson.M{"name": 1, "imageUrl": 1}},
					},
					"as": "user",
				}},
				{"$unwind": bson.M{"path": "$user", "preserveNullAndEmptyArrays": true}},
				{"$project": bson.M{"_id": 1, "value": 1, "createdAt": 1, "userId": 1, "user": 1}},
			},
			"as": "comments",
		}},
		// Ensure all fields are projected
		{"$project": bson.M{
			"_id":          1,
			"creator":      1,
			"title":        1,
			"message":      1,
			"name":         1,
			"selectedFile": 1,
			"likes":        1,
			"createdAt":    1,
			"comments":     1,
		}},
	}

	cursor, err := postCol.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "aggregation failed", "details": err.Error()})
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read aggregation results", "details": err.Error()})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"post": results[0]})
}

// like Post
// @Summary like or unkike a post
// @Description Like or un like a post  by it's id
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Success 200 {object} models.PostModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id}/likePost [patch]
func LikePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	var UserSchema = database.DB.Collection("users")
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var post models.PostModel
	postid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	// after getting post
	userID, errb := c.Locals("userId").(string)
	if !errb {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"deatils": "you are not authorized",
		})
	}

	// check
	if slices.Contains(post.Likes, userID) {
		i := sort.SearchStrings(post.Likes, userID)
		post.Likes = slices.Delete(post.Likes, i, i+1)
	} else {
		post.Likes = append(post.Likes, userID)
		//  START craete Notification
		objId, _ := primitive.ObjectIDFromHex(userID)
		var user models.UserModel

		// get nuser data
		userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
		if userResult.Err() != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"success": false,
				"message": "User Not found",
			})
		}

		userResult.Decode(&user)
		// Create Notification
		notification := models.Notification{
			MainUID:   post.Creator,
			TargetID:  post.ID.Hex(),
			Deatils:   user.Name + " Liked your Post",
			User:      models.User{Name: user.Name, Avatart: user.ImageUrl},
			CreatedAt: time.Now(),
		}
		res, err := NotificationSchema.InsertOne(ctx, notification)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Faild to create notification",
				"error":   err.Error(),
			})
		}

		// set the id fiald of the notficato object
		notification.ID = res.InsertedID.(primitive.ObjectID)
		// call grpc
		servergrpc.SendNotification(notification)
		// End create notfication
	}

	// update post
	updateLikelist := bson.M{"likes": post.Likes}
	_, err = PostSchema.UpdateOne(ctx, bson.M{"_id": postid}, bson.M{"$set": updateLikelist})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}
	err = PostSchema.FindOne(ctx, bson.M{"_id": postid}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"post": post,
	})
}

// Delete Post
// @Summary Delete  post by id
// @Description Delete post by post id need to prvided auth token for post craetor
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path string true "Post Id"
// @Failure 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /posts/{id} [delete]
func DeletePost(c *fiber.Ctx) error {

	var PostSchema = database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// authorization start
	var authPost models.PostModel
	primID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	PostSchema.FindOne(ctx, bson.M{"_id": primID}).Decode(&authPost)

	if authPost.Creator != c.Locals("userId").(string) {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "You Are Not authorized to delete this post.",
		})
	}

	//
	result, err := PostSchema.DeleteOne(ctx, bson.M{"_id": primID})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"deatils": err.Error(),
		})
	}

	if result.DeletedCount == 1 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Post Deleted Successfully!",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "can't Delete Post!",
		})
	}

}

// Delete Comment
// @Summary Delete a comment
// @Description Delete a comment by comment ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param postId path string true "Post ID"
// @Param commentId path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /comments/{postId}/comments/{commentId} [delete]
func DeleteComment(c *fiber.Ctx) error {
	commentCol := database.DB.Collection("comments")
	postCol := database.DB.Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Reduced timeout
	defer cancel()

	// Get comment ID from params
	commentID, err := primitive.ObjectIDFromHex(c.Params("commentId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid comment ID",
		})
	}

	// Get post ID from params
	postID, err := primitive.ObjectIDFromHex(c.Params("postId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Get current user ID
	currentUserID := c.Locals("userId").(string)
	currentUserObjID, _ := primitive.ObjectIDFromHex(currentUserID)

	// Find the comment and check ownership in one query
	var comment models.Comment
	err = commentCol.FindOne(ctx, bson.M{
		"_id":    commentID,
		"postId": postID,
	}).Decode(&comment)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	// Get post to check if user is post owner
	var post models.PostModel
	err = postCol.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	// Check authorization: comment owner or post owner
	if comment.UserID != currentUserObjID && post.Creator != currentUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this comment",
		})
	}

	// Delete the comment - use simple delete
	result, err := commentCol.DeleteOne(ctx, bson.M{"_id": commentID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete comment",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found or already deleted",
		})
	}

	// Return success immediately - don't fetch the whole post again
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Comment deleted successfully",
		"deletedCommentId": commentID.Hex(),
	})
}
