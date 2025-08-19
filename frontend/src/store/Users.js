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

        GetUserFollowersFollowing: async () => {
            const userd = JSON.parse(localStorage.getItem('profile'));
            const {data}= await api.fetchUserProfile(userd.result._id)

            var followers = data.user.followers || [];
            var following = data.user.following || [];
            
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
        UserData(state, payload) {
            state.User = payload
        },
        SetUserPosts(state, payload) {
            if (payload.append) {
                // Prevent duplicates when appending
                const existingIds = new Set(state.UserPosts.map(post => post._id))
                const newPosts = payload.posts.filter(post => !existingIds.has(post._id))
                state.UserPosts = [...state.UserPosts, ...newPosts]
            } else {
                // Replace posts for new user or reset
                state.UserPosts = payload.posts || []
            }
        },
        SetPostsPagination(state, payload) {
            state.PostsPagination = {
                currentPage: payload.currentPage || 1,
                numberOfPages: payload.numberOfPages || 0
            }
        },
        ResetUserPosts(state) {
            state.UserPosts = []
            state.PostsPagination = {
                currentPage: 1,
                numberOfPages: 0
            }
        },
        // Add mutation to update follow status
        UpdateUserFollowStatus(state, { userId, isFollowing, followersCount }) {
            if (state.User && state.User._id === userId) {
                state.User.isFollowing = isFollowing
                state.User.followersCount = followersCount
            }
        }
    },
    actions: {
        // Enhanced GetUserByID with better error handling
        async GetUserByID(context, { id, page = 1, append = false }) {
            try {
                console.log('Fetching user profile for ID:', id, 'Page:', page)

                // Validate ID - handle both undefined and empty string
                if (!id || id === 'undefined' || id.trim() === '') {
                    console.error('Invalid user ID provided:', id)
                    throw new Error('User ID is required')
                }

                const { data } = await api.fetchUserProfile(id, page);
                console.log('API Response:', data)

                // Validate response structure
                if (!data) {
                    throw new Error('No data received from API')
                }

                // Update user data only if not appending
                if (!append && data.user) {
                    context.commit('UserData', data.user)
                }

                // Update posts with better validation
                const posts = Array.isArray(data.posts) ? data.posts : []
                context.commit('SetUserPosts', {
                    posts: posts,
                    append: append
                })

                // Update pagination
                context.commit('SetPostsPagination', {
                    currentPage: data.currentPage || page,
                    numberOfPages: data.numberOfPages || 0
                })

                return {
                    user: data.user,
                    posts: posts,
                    currentPage: data.currentPage || page,
                    numberOfPages: data.numberOfPages || 0
                };

            } catch (error) {
                console.error('GetUserByID Error:', error);

                // Handle specific error cases
                if (error.response?.status === 502) {
                    console.error('Server is currently unavailable. Please try again later.')
                } else if (error.response?.status === 404) {
                    console.error('User not found')
                } else if (error.response?.status === 401) {
                    console.error('Unauthorized access')
                }

                throw error;
            }
        },



        // Reset user posts
        ResetUserPosts({ commit }) {
            commit('ResetUserPosts')
        },

        // Enhanced update user data
        async UpdateUserData(context, userData) {
            try {
                console.log('Updating user data:', userData)
                const { data } = await api.UpdateUser(userData);

                if (data?.user) {
                    context.commit('UserData', data.user)
                }

                return data;
            } catch (error) {
                console.error('UpdateUserData Error:', error);
                throw error;
            }
        },

        // Enhanced follow user with state updates
        async FollowUser({ commit }, ProfileID) {
            try {
                console.log('Following user:', ProfileID)
                const { data } = await api.following(ProfileID);

                // Update local state immediately for better UX
                if (data?.success) {
                    commit('UpdateUserFollowStatus', {
                        userId: ProfileID,
                        isFollowing: data.isFollowing,
                        followersCount: data.followersCount
                    });
                }

                return data;
            } catch (error) {
                console.error('FollowUser Error:', error);
                throw error;
            }
        },

        async GetTheUserSug(context, id) {
            try {
                console.log('Getting user suggestions for:', id)
                const { data } = await api.getSugUser(id);
                return data;
            } catch (error) {
                console.error('GetTheUserSug Error:', error);
                throw error;
            }
        }
    }
}

export default Users
