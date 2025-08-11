import * as api from '@/api/index.js'

const Users = {
    state: {
        User: null,
        UserPosts: [],
        PostsPagination: {
            currentPage: 1,
            numberOfPages: 0
        }
    },
    getters: {
        GetUser: (state) => () => {
            return state.User
        },
        GetUserPosts: (state) => {
            return state.UserPosts
        },
        GetPostsPagination: (state) => {
            return state.PostsPagination
        },
        // TODO GetUserFollowersFollowing 
        GetUserFollowersFollowing: async () => {
            const userd = JSON.parse(localStorage.getItem('profile'));
            var followers = userd.result.followers || [];
            var following = userd.result.following || [];
            
            const combinedArray = [...followers, ...following];
            const uniqueArray = Array.from(new Set(combinedArray));

            var userdata = [];
            for(const uid of uniqueArray){
                const {data } = await api.fetchUserProfile(uid);
                var user = {"_id": data.user._id, "name": data.user.name, "imageUrl": data.user.imageUrl};
                userdata.push(user)
            }
            return userdata;
        }
    },
    mutations: {
        UserData(state, payload){
            state.User = payload?.data
        },
        SetUserPosts(state, payload) {
            if (payload.append) {
                // Append new posts for pagination
                state.UserPosts = [...state.UserPosts, ...payload.posts]
            } else {
                // Replace posts for new user or reset
                state.UserPosts = payload.posts
            }
        },
        SetPostsPagination(state, payload) {
            state.PostsPagination = {
                currentPage: payload.currentPage,
                numberOfPages: payload.numberOfPages
            }
        },
        ResetUserPosts(state) {
            state.UserPosts = []
            state.PostsPagination = {
                currentPage: 1,
                numberOfPages: 0
            }
        }
    },
    actions: {
        // getuserbyid with posts
        async GetUserByID(context, { id, page = 1, append = false }) {
            try {
                const {data} = await api.fetchUserProfile(id, page);

                // Update user data only if not appending
                if (!append) {
                    context.commit('UserData', data.user)
                }

                // Update posts - data already contains posts array from your backend
                context.commit('SetUserPosts', { 
                    posts: data.posts || [], 
                    append: append 
                })

                // Update pagination - using the exact field names from your backend
                context.commit('SetPostsPagination', {
                    currentPage: data.currentPage || page,
                    numberOfPages: data.numberOfPages || 0
                })

                return data;
            } catch (error) {
                console.log(error);
                return error;
            }
        },
        // Reset user posts (for when switching profiles)
        ResetUserPosts(context) {
            context.commit('ResetUserPosts')
        },
        // update user data
        async UpdateUserData(context, userData) {
            try {
                const {data} = await api.UpdateUser(userData);

                context.commit('UserData', data.user)

                return data;
            } catch (error) {
                console.log(error);
                return error;
            }
        },
        // following user
        async FollowUser(context, ProfileID) {
            try {
                const {data} = await api.following(ProfileID )

                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async GetTheUserSug(context, id){
            try {
                const {data} = await api.getSugUser(id)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        }
    }
}

export default Users


// import * as api from '@/api/index.js'

// const Users = {
//     state: {User: null},
//     getters: {
//         GetUser: (state) => () => {
//             return state.User
//         },
//         // TODO GetUserFollowersFollowing 
//         GetUserFollowersFollowing: async () => {
//             const userd = JSON.parse(localStorage.getItem('profile'));
//             var followers = userd.result.followers || [];
//             var following = userd.result.following || [];
            
//             const combinedArray = [...followers, ...following];
//             const uniqueArray = Array.from(new Set(combinedArray));

//             var userdata = [];
//             for(const uid of uniqueArray){
//                 const {data } = await api.fetchUserProfile(uid);
//                 var user = {"_id": data.user._id, "name": data.user.name, "imageUrl": data.user.imageUrl};
//                 userdata.push(user)
//             }
//             return userdata;
//         }
//     },
//     mutations: {
//         UserData(state, payload){
//             state.User = payload?.data
//         }
//     },
//     actions: {
//         // getuserbyid
//         async GetUserByID(context, id) {
//             try {
//                 const {data} = await api.fetchUserProfile(id);

//                 context.commit('UserData', data.user)

//                 return data;
//             } catch (error) {
//                 console.log(error);
//                 return error;
//             }
//         },
//         // update user data
//         async UpdateUserData(context, userData) {
//             try {
//                 const {data} = await api.UpdateUser(userData);

//                 context.commit('UserData', data.user)

//                 return data;
//             } catch (error) {
//                 console.log(error);
//                 return error;
//             }
//         },
//         // following user
//         async FollowUser(context, ProfileID) {
//             try {
//                 const {data} = await api.following(ProfileID )

//                 return data
//             } catch (error) {
//                 console.log(error)
//                 return error
//             }
//         },
//         async GetTheUserSug(context, id){
//             try {
//                 const {data} = await api.getSugUser(id)
//                 return data
//             } catch (error) {
//                 console.log(error)
//                 return error
//             }
//         }
//     }
// }



// export default Users



