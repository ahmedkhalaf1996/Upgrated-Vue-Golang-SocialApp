<template>
  <q-page class="constrain q-pa-md">
    <div class="row q-col-gutter-lg constrain">
      <!-- Error State -->
      <div v-if="error" class="col-12 q-pa-lg text-center text-negative">
        <q-icon name="eva-alert-circle-outline" size="48px" />
        <div class="q-mt-md text-h6">{{ error }}</div>
        <q-btn 
          @click="resetAndLoadProfile" 
          color="primary" 
          outline 
          class="q-mt-md"
        >
          Try Again
        </q-btn>
      </div>

      <!-- Profile Content -->
      <template v-else>
        <ShowProfile
          :userData="userData"
          :userPosts="userPosts"
          :isSameUser="isSameUser"
          @EditProfile="EditMode = !EditMode"
          @update-user="updateUserLocal"
          v-if="!EditMode && userData"
        />
        <EditProfile
          :userData="userData"
          :isSameUser="isSameUser"
          @EditProfile="EditMode = !EditMode"
          @update-user="updateUserLocal"
          v-else-if="EditMode && userData"
        />
        
        <div class="col-12" v-if="userData">
          <q-separator inset />
        </div>
        
        <!-- Loading skeleton for initial posts -->
        <div v-if="!load && !error" class="col-12">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-sm-6 col-md-4" v-for="i in 3" :key="i">
              <q-card>
                <q-skeleton height="200px" square />
                <q-card-section>
                  <q-skeleton type="text" class="text-h6" />
                  <q-skeleton type="text" width="50%" class="text-subtitle2" />
                </q-card-section>
              </q-card>
            </div>
          </div>
        </div>
        
        <!-- Posts grid -->
        <div
          v-else-if="load && userPosts.length > 0"
          class="col-12 col-sm-6 col-md-4"
          v-for="post in userPosts"
          :key="post._id"
        >
          <Post :post="post" />
        </div>
        
        <!-- Loading indicator for more posts -->
        <div v-if="loadingMore" class="col-12 q-pa-lg text-center">
          <q-spinner-hourglass color="primary" size="3em" />
          <div class="q-mt-md text-grey-7">
            Loading more posts...
          </div>
        </div>
        
        <!-- End of posts indicator -->
        <div v-if="hasReachedEnd && userPosts.length > 0 && load" class="col-12 q-pa-md text-center text-grey-6">
          <q-icon name="eva-inbox-outline" size="24px" />
          <div class="q-mt-sm">No more posts</div>
        </div>
        
        <!-- No posts message -->
        <div v-if="load && userPosts.length === 0 && userData && !error" class="col-12 q-pa-lg text-center text-grey-6">
          <q-icon name="eva-image-outline" size="48px" />
          <div class="q-mt-md text-h6">No posts yet</div>
          <div class="text-body2">
            {{ isSameUser ? "You haven't" : "This user hasn't" }} shared any posts.
          </div>
        </div>
      </template>
    </div>
  </q-page>
</template>

<script>
import { mapGetters, mapMutations, mapActions } from 'vuex';
import Post from '@/components/post/Post.vue';
import ShowProfile from '@/components/user/ShowProfile.vue'
import EditProfile from '@/components/user/EditProfile.vue';

export default {
  name: 'ProfileView',
  data(){
    return {
      userData: null,
      isSameUser: false,
      EditMode: false,
      load: false,
      loadingMore: false,
      currentPage: 1,
      maxPages: 0,
      hasReachedEnd: false,
      error: null,
      scrollListener: null,
    }
  },
  watch:{
    $route: {
      handler(newRoute, oldRoute) {
        // Only reset if the user ID actually changed
        if (newRoute.params.id !== oldRoute?.params.id) {
          this.resetAndLoadProfile()
        }
      },
      immediate: true
    }
  },
  async mounted(){
    // console.log("Profile mounted for userid:", this.$route.params.id)
    this.SetData();
    await this.resetAndLoadProfile()
  },
  computed: {
    ...mapGetters(['GetUserData', 'GetUserPosts', 'GetPostsPagination']),
    userPosts() {
      return this.GetUserPosts || [];
    }
  },
  methods: {
    ...mapMutations(['SetData']),
    ...mapActions(['GetUserByID', 'ResetUserPosts']),
    
    // Reset and load profile data
    async resetAndLoadProfile() {
      try {
        // Reset state
        this.error = null;
        this.ResetUserPosts();
        this.load = false;
        this.loadingMore = false;
        this.currentPage = 1;
        this.hasReachedEnd = false;
        this.userData = null;
        
        // Remove existing scroll listener
        this.removeScrollListener();
        
        // Load initial data
        await this.GetAll(false);
        
        // Add scroll listener after successful load
        if (!this.error) {
          this.addScrollListener();
        }
        
      } catch (error) {
        console.error('Failed to reset and load profile:', error);
        this.error = 'Failed to load profile. Please try again.';
      }
    },
    
    // Get All User data & posts with enhanced error handling
    async GetAll(append = false){
      try {
        const LoggedUserID = this.GetUserData()?.result?._id;
        const profileId = this.$route.params.id;
        
        // console.log('GetAll - LoggedUser:', LoggedUserID, 'ProfileID:', profileId);
        
        // Validate profileId more thoroughly
        if (!profileId || profileId === 'undefined' || profileId.trim() === '') {
          console.error('Invalid profile ID in route:', profileId)
          throw new Error('Profile ID is missing or invalid');
        }
        
        const data = await this.GetUserByID({
          id: profileId,
          page: this.currentPage,
          append: append
        });
        
        // console.log('GetAll - API Response:', data);
        
        // Update local component data
        if (!append && data?.user) {
          this.userData = data.user;
          this.isSameUser = String(LoggedUserID) === String(profileId);
          console.log('Is same user:', this.isSameUser);
        }
        
        // Update pagination data
        this.maxPages = data?.numberOfPages || 0;
        this.hasReachedEnd = this.currentPage >= this.maxPages;
        
        this.load = true;
        this.error = null;
        
        return data;
        
      } catch (error) {
        console.error("Error in GetAll:", error);
        
        // Set appropriate error message based on error type
        if (error.response?.status === 502) {
          this.error = 'Server is currently unavailable. Please try again later.';
        } else if (error.response?.status === 404) {
          this.error = 'User not found.';
        } else if (error.response?.status === 401) {
          this.error = 'You are not authorized to view this profile.';
        } else if (error.message?.includes('Network Error')) {
          this.error = 'Network connection error. Please check your internet connection.';
        } else {
          this.error = 'Failed to load profile data. Please try again.';
        }
        
        this.load = true; // Set load to true to show error state
        throw error;
      }
    },
    
    async loadMorePosts() {
      // Prevent multiple requests
      if (this.loadingMore || this.hasReachedEnd || this.error) {
        return;
      }
      
      // Check if there are more pages
      if (this.currentPage < this.maxPages) {
        this.loadingMore = true;
        this.currentPage++;
        
        try {
          await this.GetAll(true);
          // Optional delay for better UX
          await new Promise(resolve => setTimeout(resolve, 500));
        } catch (error) {
          console.error('Error loading more posts:', error);
          this.currentPage--; // Revert on error
          // Don't show error for pagination failures, just log them
        } finally {
          this.loadingMore = false;
        }
      }
    },
    
    handleScroll() {
      const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
      const windowHeight = window.innerHeight;
      const documentHeight = document.documentElement.scrollHeight;
      
      // Load more when user is 300px from bottom
      if (scrollTop + windowHeight >= documentHeight - 300) {
        this.loadMorePosts();
      }
    },

    addScrollListener() {
      if (!this.scrollListener) {
        this.scrollListener = this.handleScroll.bind(this);
        window.addEventListener('scroll', this.scrollListener);
      }
    },

    removeScrollListener() {
      if (this.scrollListener) {
        window.removeEventListener('scroll', this.scrollListener);
        this.scrollListener = null;
      }
    },
    
    updateUserLocal(updatedData){
      if (updatedData?.data) {
        this.userData = updatedData.data;
      }
    }
  },
  
  beforeUnmount() {
    // Remove scroll listener
    this.removeScrollListener();
  },
  
  components: { ShowProfile, EditProfile, Post }
}
</script>

<style scoped>
.q-page {
  scroll-behavior: smooth;
}

.text-negative {
  color: #c10015;
}
</style>