<template>
   <q-page class="constrain q-pa-md resize-observer-fix">
     <div class=" row q-col-gutter-lg">
      <div class="col-3">
        <SideBar />
      </div>
      <div v-if="!load" class="col-6  q-mx-auto">
        <div class="q-pa-md">
          <q-card>
            <q-item>
              <q-item-section avatar>
                <q-skeleton type="QAvatar"/>
              </q-item-section>
              <q-item-section>
                <q-item-label>
                  <q-skeleton type="text" />
                </q-item-label>
                <q-item-label caption>
                  <q-skeleton type="text" />
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-skeleton height="200px" square />
            <q-card-actions class="q-gutter-md">
              <q-skeleton type="QBtn"/>
              <q-skeleton type="QBtn"/>
            </q-card-actions>
          </q-card>
        </div>
      </div>
      <div v-else class="col-6 q-mx-auto">
        <Post v-for="post in posts" :key="post._id" :post="post" />
        
        <!-- Loading indicator for more posts -->
        <div v-if="loadingMore" class="q-pa-lg text-center">
          <q-spinner-hourglass color="primary" size="3em" />
          <div class="q-mt-md text-grey-7">
            Loading more posts...
          </div>
        </div>
        
        <!-- End of posts indicator -->
        <div v-if="hasReachedEnd && posts.length > 0" class="q-pa-md text-center text-grey-6">
          <q-icon name="eva-inbox-outline" size="24px" />
          <div class="q-mt-sm">No more posts</div>
        </div>
      </div>
      <div class="col-3">
        <Rightbar />
      </div> 
    </div>
     <div class="q-pa-lg flex justify-center fixed-bottom">
      <Add @Created="onPostCreated"/>
     </div>
   </q-page>
</template>

<script>
import Add from '@/components/post/Add.vue'
import Post from '@/components/post/Post.vue'
import SideBar from '@/components/sideBar/SideBar.vue';
import Rightbar from '@/components/rightbar/Rightbar.vue';
import { mapActions } from 'vuex';

export default {
  name: 'HomeView',
  data(){
    return {
      currentPage: 1,
      maxPages: 0,
      posts: [],
      load: false,
      loadingMore: false,
      hasReachedEnd: false
    }
  },
  components: {
    Add,
    Post,
    SideBar,
    Rightbar,
  },
  methods:{
    ...mapActions(['getPosts']),
    
    async GetAllPosts(append = false){
      console.log("Get All Posts Called", "Page:", this.currentPage, "Append:", append)
      
      try {
        const data = await this.getPosts(this.currentPage)
        console.log("post data", data)
        
        if(data?.data){
          this.maxPages = data?.numberOfPages;
          
          if (append) {
            // Append new posts to existing ones
            this.posts = [...this.posts, ...data.data];
          } else {
            // Replace posts (initial load)
            this.posts = data.data; 
          }
          
          // Check if we've reached the end
          this.hasReachedEnd = this.currentPage >= this.maxPages;
        }

        if(data){
          this.load = true;
        }
      } catch (error) {
        console.error("Error loading posts:", error);
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
          await this.GetAllPosts(true);
          // Wait 2 seconds
          await new Promise(resolve => setTimeout(resolve, 2000));
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
    
    onPostCreated() {
      // Reset and reload
      this.currentPage = 1;
      this.hasReachedEnd = false;
      this.posts = [];
      this.GetAllPosts(false);
    }
  },
  
  async mounted(){
    // Load initial posts
    setTimeout(async () => {
      await this.GetAllPosts();
      
      // Add scroll listener
      window.addEventListener('scroll', this.handleScroll);
    }, 5000);
  },
  
  beforeUnmount() {
    // Remove scroll listener
    window.removeEventListener('scroll', this.handleScroll);
  }
}
</script>

<style scoped>
.q-page {
  scroll-behavior: smooth;
}
</style>
