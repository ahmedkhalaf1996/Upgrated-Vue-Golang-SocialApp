<template>
  <q-page class="constrain q-pa-md">
    <div class="row q-col-gutter-lg constrain">
      <ShowProfile
        :userData="userData"
        :userPosts="userPosts"
        :isSameUser="isSameUser"
        @EditProfile="EditMode = !EditMode"
        @update-user="updateUserLocal"
        v-if="!EditMode"
      />
      <EditProfile
        :userData="userData"
        :isSameUser="isSameUser"
        @EditProfile="EditMode = !EditMode"
        @update-user="updateUserLocal"
        v-else
      />
      <div class="col-12">
        <q-separator inset />
      </div>
      
      <!-- Loading skeleton for initial posts -->
      <div v-if="!load" class="col-12">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-sm-6 col-md-4" v-for="i in 2" :key="i">
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
        v-else
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
      <div v-if="hasReachedEnd && userPosts.length > 0" class="col-12 q-pa-md text-center text-grey-6">
        <q-icon name="eva-inbox-outline" size="24px" />
        <div class="q-mt-sm">No more posts</div>
      </div>
      
      <!-- No posts message -->
      <div v-if="load && userPosts.length === 0" class="col-12 q-pa-lg text-center text-grey-6">
        <q-icon name="eva-image-outline" size="48px" />
        <div class="q-mt-md text-h6">No posts yet</div>
        <div class="text-body2">This user hasn't shared any posts.</div>
      </div>
    </div>
  </q-page>
</template>

<script>
// @ is an alias to /src
import { mapGetters, mapMutations, mapActions } from 'vuex';
import Post from '@/components/post/Post.vue';
import ShowProfile from '@/components/user/ShowProfile.vue'
import EditProfile from '@/components/user/EditProfile.vue';

export default {
  name: 'ProfileView',
  data(){
    return {
      userData: [],
      isSameUser: false,
      EditMode: false,
      load: false,
      loadingMore: false,
      currentPage: 1,
      maxPages: 0,
      hasReachedEnd: false,
    }
  },
  watch:{
    $route(){
      this.resetAndLoadProfile()
    }
  },
  mounted(){
    console.log("userid", this.$route.params.id)
    this.SetData();
    this.resetAndLoadProfile()
  },
  computed: {
    ...mapGetters(['GetUserData', 'GetUserPosts', 'GetPostsPagination']),
    userPosts() {
      return this.GetUserPosts;
    }
  },
  methods: {
    ...mapMutations(['SetData']),
    ...mapActions(['GetUserByID', 'ResetUserPosts']),
    
    // Reset and load profile data
    async resetAndLoadProfile() {
      // Reset pagination state
      this.ResetUserPosts();
      this.load = false;
      this.loadingMore = false;
      this.currentPage = 1;
      this.hasReachedEnd = false;
      
      // Load initial data
      await this.GetAll(false);
      
      // Add scroll listener after initial load
      setTimeout(() => {
        window.addEventListener('scroll', this.handleScroll);
      }, 1000);
    },
    
    // Get All User data & posts
    async GetAll(append = false){
      const LogedUserID = this.GetUserData()?.result?._id
      console.log('LogedUser', LogedUserID)
      
      const profileid = this.$route.params.id
      
      try {
        const data = await this.GetUserByID({
          id: profileid,
          page: this.currentPage,
          append: append
        });
        
        // Update local component data
        if (!append) {
          this.userData = data?.user;
        }
        
        // Update pagination data
        this.maxPages = data?.numberOfPages || 0;
        this.hasReachedEnd = this.currentPage >= this.maxPages;
        
        this.isSameUser = String(LogedUserID) == String(profileid);
        console.log('is same', this.isSameUser)
        
        this.load = true;
        
      } catch (error) {
        console.error("Error loading profile:", error);
      }
    },
    
    async loadMorePosts() {
      // Prevent multiple requests
      if (this.loadingMore || this.hasReachedEnd) {
        return;
      }
      
      // Check if there are more pages
      if (this.currentPage < this.maxPages) {
        this.loadingMore = true;
        this.currentPage++;
        
        try {
          await this.GetAll(true);
          // Optional delay for better UX
          await new Promise(resolve => setTimeout(resolve, 1000));
        } catch (error) {
          console.error('Error loading more posts:', error);
          this.currentPage--; // Revert on error
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
    
    updateUserLocal(updatedData){
      this.userData = updatedData.data
    }
  },
  
  beforeUnmount() {
    // Remove scroll listener
    window.removeEventListener('scroll', this.handleScroll);
  },
  
  components: { ShowProfile, EditProfile, Post }
}
</script>

<style scoped>
.q-page {
  scroll-behavior: smooth;
}
</style>




<!-- <template>
  <q-page class="constrain q-pa-md ">
    <div class="row q-col-gutter-lg constrain">
      <ShowProfile
        :userData="userData"
        :userPosts="userPosts"
        :isSameUser="isSameUser"
        @EditProfile="EditMode = !EditMode"
        @update-user="updateUserLocal"
        v-if="!EditMode"
      />
      <EditProfile
        :userData="userData"
        :isSameUser="isSameUser"
        @EditProfile="EditMode = !EditMode"
        @update-user="updateUserLocal"
        v-else
      />
      <div class="col-12">
        <q-separator inset />
      </div>
      <div 
        class="col-12 col-sm-6 col-md-4" 
        v-for="post in userPosts" 
        :key="post._id"
      >
        <Post :post="post" />
      </div>
    </div>
  </q-page>
</template>


  <script>
  // @ is an alias to /src
  import { mapGetters, mapMutations, mapActions } from 'vuex';
  import Post from '@/components/post/Post.vue';
  import ShowProfile from '@/components/user/ShowProfile.vue'
  import EditProfile from '@/components/user/EditProfile.vue';
  export default {
    name: 'ProfileView',
    data(){
      return {
        userPosts:[],
        userData:[],
        isSameUser: false,
        EditMode: false,
      }
    },
    watch:{
      $route(){
        this.GetAll()
      }
    },
    mounted(){
        console.log("userid", this.$route.params.id)
        this.SetData();
        this.GetAll()
    },
    created(){
      this.GetAll()
    },
    computed: {
      ...mapGetters(['GetUserData'])
    },
    methods: {
      ...mapMutations(['SetData']),
      ...mapActions(['GetUserByID']),
      // Get All User data & posts
      async GetAll(){
        const LogedUserID = this.GetUserData()?.result?._id
        console.log('LogedUser', LogedUserID)

        const profileid =  this.$route.params.id

        const data = await this.GetUserByID(profileid)

        this.userData = data?.user 
        this.userPosts = data?.posts

        this.isSameUser = String(LogedUserID) == String(profileid);

        console.log('is same', this.isSameUser)
      },
      updateUserLocal(updatedData){
        this.userData = updatedData.data
      }
    },
    components:{ShowProfile, EditProfile, Post}
  }
  </script> -->