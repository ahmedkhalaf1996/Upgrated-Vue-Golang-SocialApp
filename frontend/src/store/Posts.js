// Updated store/post.js

import * as api from '../api/index.js';

const Posts = {
    state: {
        isLoading: true, 
        post: [], 
        posts: [], 
        SearchResult: []
    },
    getters: {
        GetPost: (state) => () => {
            return {...state.post};
        },
        GetAllPosts: (state) => () => {
            return {...state.posts};
        },
        GetSearchData: (state) => () => {
            return {...state.SearchResult};
        },
    },
    mutations: {
        Post(state, payload){
            state.post = payload;
        },
        Posts(state, payload){
            state.posts = payload;
        },
        Search(state, payload){
            state.SearchResult = payload;
        }
    },
    actions: {
        async getPost(context, id){
            try {
                let {data} = await api.fetchPost(id);
                // Backend returns {post: postData}, so we need to extract the post
                const postData = data.post || data;
                context.commit('Post', postData);
                return postData;
            } catch (error) {
                console.log(error)
                throw error;
            }
        },
        
        async getPosts(context, page){
            try {
                const user = JSON.parse(localStorage.getItem('profile'));
                const userId = user?.result?._id;

                if (userId){
                    const {data} = await api.fetchPosts(page, userId);
                    context.commit('Posts', data)
                    return data;
                }
            } catch (error) {
                console.log(error)
                throw error;
            }
        },
        
        async getPostsUsersBySearch(context, searchData){
            try {
                const {data} = await api.fetchPostsUsersbySearch(searchData);
                context.commit('Search', data)
                return data;
            } catch (error) {
                console.log(error)
                throw error;
            }
        },
        
        async createPost(context, post){
            try {
                const {data} = await api.createPost(post);
                context.commit('Post', data);
                return data;
            } catch (error) {
                console.log(error);
                throw error;
            }
        },
        
        async updatePost(context, Data) {
            try {
                const user = JSON.parse(localStorage.getItem('profile'));
                const userId = user?.result?._id;

                const PostData = {
                    "title": Data.title,
                    "message": Data.message,
                    "creator": userId,
                    "selectedFile": Data.selectedFile,
                }

                const {data} = await api.updatePost(Data.id, PostData);
                context.commit('Post', data);
                return data;
            } catch (error) {
                console.log(error);
                throw error;
            }
        },
        
        async LikePostByUser(context, id){
            try {
                const {data} = await api.likePost(id);
                // Backend returns {post: postData}, so we need to extract the post
                const postData = data.post || data;
                context.commit('Post', postData);
                return postData;
            } catch (error) {
                console.log(error);
                throw error;
            }
        },
        
        async commentPost(context, form) {
            try {
                const {data} = await api.comment(form.value, form.id);
                // Backend returns {post: postData}, so we need to extract the post
                const postData = data.post || data;
                context.commit('Post', postData);
                return postData;
            } catch (error) {
                console.log(error);
                throw error;
            }
        },
        
        async deleteComment(context, {postId, commentId}) {
            try {
                console.log('Deleting comment:', { postId, commentId }); // Debug log
                const {data} = await api.deleteComment(postId, commentId);
                console.log('Delete response:', data); // Debug log
                return data;
            } catch (error) {
                console.error('Store deleteComment error:', error);
                // Log more details about the error
                if (error.response) {
                    console.error('Error response:', error.response.data);
                    console.error('Error status:', error.response.status);
                }
                throw error;
            }
        },
        
        async deltePost(context, id){
            try {
                await api.deltePost(id);
            } catch (error) {
                console.log(error);
                throw error;
            }
        }
    }
}

export default Posts;



// import * as api from '../api/index.js';

// const Posts = {
//     state: {isLoading: true, post:[], posts: [], SearchResult:[]},
//     getters: {
//         GetPost: (state) => () => {
//             return {...state.post};
//         },
//         GetAllPosts: (state) => () => {
//             return {...state.posts};
//         },
//         GetSearchData: (state) => () => {
//             return {...state.SearchResult};
//         },
//     },
//     mutations: {
//         Post(state, payload){
//             state.post = payload;
//         },
//         Posts(state, payload){
//             state.Posts = payload;
//         },
//         Search(state, payload){
//             state.SearchResult = payload;
//         }
//     },
//     actions: {
//         async getPost(context, id){
//             try {
//                 let {data} = await api.fetchPost(id);
//                 context.commit('Post',data)

//                 return data
//             } catch (error) {
//                 console.log(error)
//             }
//         },
//         async getPosts(context, page){
//             try {
//                 const user = JSON.parse(localStorage.getItem('profile'));
//                 const userId = user?.result?._id;

//                 if (userId){
//                     const {data}  = await api.fetchPosts(page, userId);
//                     context.commit('Posts', data)
//                     return data;
//                 }
//             } catch (error) {
//                 console.log(error)
//             }
//         },
//         // this for users & posts
//         async getPostsUsersBySearch(context, searchData){
//             try {
//                 const {data} = await api.fetchPostsUsersbySearch(searchData);
//                 // console.log(data)
//                 context.commit('Search', data)
//                 return data;
//             } catch (error) {
//                 console.log(error)
//             }
//         },
//         async createPost(context, post){
//             try {
//                 const {data} = await api.createPost(post);
//                 context.commit('Post', data);
//                 return data;
//             } catch (error) {
//                 return error;
//             }
//         },
//         async updatePost(context, Data) {
//             const user = JSON.parse(localStorage.getItem('profile'));
//             const userId = user?.result?._id;

//             const PostData = {
//                 "title": Data.title,
//                 "message": Data.message,
//                 "creator": userId,
//                 "selectedFile": Data.selectedFile,
//             }

//             const post = await api.updatePost(Data.id, PostData);
//             context.commit('Post', post);

//             return post;

//         },
//         async LikePostByUser(context, id){
//             // const user = JSON.parse(localStorage.getItem('profile'));
//             const {data} = await api.likePost(id)
//             context.commit('Post', data);
//             console.log('data', data)
//         },
//         async commentPost(context, form) {
//             const { data } = await api.comment(form.Value, form.id)
//             context.commit('Post', data);
//         },
//         async deltePost(context, id){
//             await api.deltePost(id);
//         }
//     }


// }

// export default Posts;


