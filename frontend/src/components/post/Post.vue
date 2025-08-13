<template>
    <div>
         <q-card v-if="!EditPost" class="card-post q-mb-md" flat bordered>
            <q-item>
                <q-item-section avatar>
                    <q-avatar>
                        <!-- Use user data from populated field -->
                        <img v-if="localPost?.user?.imageUrl" :src="localPost.user.imageUrl" />
                        <img v-else src="https://cdn-icons-png.flaticon.com/512/1077/1077063.png" />
                    </q-avatar>
                </q-item-section>

                <q-item-section>
                    <!-- Use populated user name -->
                    <q-item-label class="text-bold">{{ localPost?.user?.name || 'Loading...' }}</q-item-label>
                    <q-item-label caption>
                        {{  getTime()  }}
                    </q-item-label>
                </q-item-section>
            </q-item>

            <q-separator />
            <q-img 
                v-if="localPost?.selectedFile" 
                style="cursor: pointer;" 
                @click="GoToDeatils" 
                :src="localPost.selectedFile" 
            />

            <q-card-section>
                <div class="text-h6">{{ localPost?.title || '' }}</div>
                <div class="text-subtitle1">{{ localPost?.message || '' }}</div>
                
                <!-- Improved Actions Section -->
                <div class="row items-center q-mt-md q-mb-md">
                    <q-btn v-if="!UserLike" @click="Like" flat round color="red" icon="eva-heart-outline" size="sm">
                        {{ LikesCount() }}
                    </q-btn>
                    <q-btn v-else @click="Like" flat round color="red" icon="eva-heart" size="sm">
                        {{ LikesCount() }}
                    </q-btn>
                </div>
                
                <q-separator class="q-my-md" />
                
                <div class="comments-section">
                    <div class="text-h6 q-mb-md">Comments</div>
                    
                    <!-- Show comments based on showAllComments state -->
                    <div 
                        v-for="comment in displayedComments"
                        :key="comment._id"
                        class="comment-item q-mb-md"
                    >
                        <q-item>
                            <q-item-section avatar>
                                <q-avatar size="32px">
                                    <img v-if="comment.user?.imageUrl" :src="comment.user.imageUrl" />
                                    <img v-else src="https://cdn-icons-png.flaticon.com/512/1077/1077063.png" />
                                </q-avatar>
                            </q-item-section>
                            <q-item-section>
                                <q-item-label class="text-bold text-caption">{{ comment.user?.name || 'Anonymous' }}</q-item-label>
                                <q-item-label class="text-body2">{{ comment.value || '' }}</q-item-label>
                                <q-item-label caption>
                                    {{ getCommentTime(comment.createdAt) }}
                                </q-item-label>
                            </q-item-section>
                            <q-item-section side v-if="canDeleteComment(comment)">
                                <q-btn 
                                    flat 
                                    round 
                                    dense 
                                    color="negative" 
                                    icon="delete" 
                                    size="sm"
                                    @click="deleteCommentConfirm(comment._id)"
                                    :loading="deletingComments[comment._id]"
                                />
                            </q-item-section>
                        </q-item>
                        <q-separator />
                    </div>
                    
                    <!-- Show more/less comments button -->
                    <div v-if="hasMoreComments" class="text-center q-mb-md">
                        <q-btn 
                            v-if="!showAllComments"
                            flat 
                            color="primary" 
                            @click="showAllComments = true"
                            class="text-caption"
                        >
                            Show {{ remainingCommentsCount }} more comments
                        </q-btn>
                        <q-btn 
                            v-else
                            flat 
                            color="primary" 
                            @click="showAllComments = false"
                            class="text-caption"
                        >
                            Show less
                        </q-btn>
                    </div>
                    
                    <div v-if="!localPost.comments || localPost.comments.length === 0" class="text-grey-6 text-center q-py-md">
                        No comments yet. Be the first to comment!
                    </div>
                </div>
            </q-card-section>

            <q-card-section>
                <q-input 
                    outlined 
                    v-model="form.text" 
                    label="Add a comment..." 
                    :loading="commentLoading"
                    @keyup.enter="AddComment"
                >
                    <template v-slot:append>
                        <q-btn 
                            v-if="form.text.trim() !== ''" 
                            @click="AddComment" 
                            flat 
                            round 
                            color="primary"
                            icon="send"
                            :loading="commentLoading"
                        />
                    </template>
                </q-input>
            </q-card-section>

         </q-card>
          <div v-else class="q-pa-md items-start q-gutter-md">
             <q-card class="my-card col-12">
                <q-card-section>
                    <div class="text-h6">Edit Post</div>
                    <q-input dense v-model="localPost.title" autofocus placeholder="Post Title" /> 
                    <div>
                        <q-input v-model="localPost.message"
                             placeholder="What's on your mind!"
                             type="textarea"
                             />
                    </div>
                    <div class="q-pa-md">
                        <q-file 
                        v-model="file"
                        label="Pick Image"
                        filled 
                        />
                    </div>
                          
                    <div v-if="localPost.selectedFile">
                        <q-img 
                        :src="localPost.selectedFile"
                        spinner-color="red"
                        style="height: 140px; max-width: 150px;"
                        />
                    </div>

                    <q-btn flat label="Update" v-close-popup @click="FireUpdate" />
                </q-card-section>
             </q-card>
          </div>
    </div>
</template>

<script>
import moment from 'moment';
import { mapActions, mapGetters } from 'vuex';

export default {
    name:'PostComponent',
    props:['post', 'EditPost'],
    data(){
        return {
            form:{text:''},
            file:null,
            UserLike:false,
            localPost: {},
            commentLoading: false,
            deletingComments: {},
            showAllComments: false,
            isLoading: true
        }
    },
    watch:{
        file(){
            this.ConvertToBase64()
        },
        post: {
            handler(newPost) {
                if (newPost) {
                    this.localPost = JSON.parse(JSON.stringify(newPost));
                }
            },
            deep: true,
            immediate: true
        }
    },

    methods:{
        ...mapActions(['GetUserDataOnly', 'LikePostByUser', 'commentPost', 'updatePost', 'deleteComment']),

        GoToDeatils(){
            this.$router.push({path:`/PostDeatils/${ this.localPost?._id}`})
        },
        
        async FireUpdate(){
            const PostData = {
                id:  this.localPost._id,
                title:  this.localPost.title,
                selectedFile:  this.localPost.selectedFile,
                message:  this.localPost.message,
            }

            const res = await this.updatePost(PostData)
            if(res){
                this.$emit('changeEdit')
            }
        },
        
        getTime(){
            return this.localPost?.createdAt ? moment(this.localPost.createdAt).fromNow() : 'Just now'
        },
        
        getCommentTime(createdAt){
            return createdAt ? moment(createdAt).fromNow() : 'Just now'
        },
        
        Like(){
            if (!this.GetUserData()?.result?._id) return;
            
            this.LikePostByUser(this.localPost._id);
            const uid = this.GetUserData().result._id;
            if(this.UserLike){
                 this.localPost.likes = (this.localPost.likes || []).filter(id => id != uid)
           } else {
             this.localPost.likes = this.localPost.likes || [];
             this.localPost.likes.push(uid)
           }
           this.UserLike = !this.UserLike
        },
        
        LikesCount(){
            if(this.localPost.likes?.length > 0){
                return String(this.localPost.likes?.length)
            }
            return '0'
        },
        
        async AddComment(){
            if (!this.form.text.trim()) return;
            
            this.commentLoading = true;
            try {
                const response = await this.commentPost({
                    value: this.form.text.trim(), 
                    id: this.localPost._id
                });
                
                if (response) {
                    this.localPost = response;
                    this.$emit('postUpdated', response);
                }
                
                this.form.text = '';
                
                this.$q.notify({
                    color: 'positive',
                    message: 'Comment added successfully',
                });
            } catch (error) {
                console.error('Error adding comment:', error);
                this.$q.notify({
                    color: 'negative',
                    message: 'Failed to add comment',
                    icon: 'error'
                });
            } finally {
                this.commentLoading = false;
            }
        },

        canDeleteComment(comment) {
            const currentUserId = this.GetUserData()?.result?._id;
            if (!currentUserId) return false;
            // Compare with string version of creator ObjectID
            return comment.userId === currentUserId || String(this.localPost.creator) === currentUserId;
        },

        deleteCommentConfirm(commentId) {
            if (confirm('Are you sure you want to delete this comment?')) {
                this.deleteCommentNow(commentId);
            }
        },

        async deleteCommentNow(commentId) {
            this.deletingComments[commentId] = true;
            
            try {
                 await this.deleteComment({
                    postId: this.localPost._id,
                    commentId: commentId
                });
                
                this.localPost.comments = (this.localPost.comments || []).filter(
                    comment => comment._id !== commentId
                );
                
                this.$q.notify({
                    color: 'positive',
                    message: 'Comment deleted successfully',
                });
            } catch (error) {
                console.error('Error deleting comment:', error);
                this.$q.notify({
                    color: 'negative',
                    message: 'Failed to delete comment: ' + (error.response?.data?.error || error.message),
                    icon: 'error'
                });
            } finally {
                delete this.deletingComments[commentId];
                this.$forceUpdate();
            }
        },

        ConvertToBase64(){
            var reader = [];
            reader = new FileReader();
            reader.readAsDataURL(this.file);

            new Promise(()=> {
                reader.onload = ()=> {
                     this.localPost.selectedFile = reader.result
                }
            })
        }
    },
    computed:{
        ...mapGetters(['GetUserData']),
        
        displayedComments() {
            if (!this.localPost.comments) return [];
            if (this.showAllComments || this.localPost.comments.length <= 2) {
                return this.localPost.comments;
            }
            return this.localPost.comments.slice(0, 2);
        },
        
        hasMoreComments() {
            return this.localPost.comments && this.localPost.comments.length > 2;
        },
        
        remainingCommentsCount() {
            if (!this.localPost.comments) return 0;
            return this.localPost.comments.length - 2;
        }
    }, 
    async mounted(){
        if (this.post) {
            this.localPost = JSON.parse(JSON.stringify(this.post));
            
            // Check if user liked the post
            const uid = this.GetUserData()?.result?._id;
            if (uid && this.localPost.likes) {
                var isLike = this.localPost.likes.find((like)=> like == uid);
                this.UserLike = !!isLike;
            }
            
            this.isLoading = false;
        }
    }
}
</script>

<style scoped>
.comment-item {
    background-color: #f9f9f9;
    border-radius: 8px;
    padding: 4px;
}

.comments-section {
    max-height: 400px;
    overflow-y: auto;
}

.card-post {
    border-radius: 8px;
}
</style>


