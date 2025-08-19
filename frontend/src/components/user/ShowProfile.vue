<template>
  <div class="row col-12">
    <div class="col-4 text-center resize-observer-fix">
      <q-avatar size="150px">
        <img v-if="userData?.imageUrl" :src="userData?.imageUrl" loading="lazy">
        <img v-else src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png" loading="lazy">
      </q-avatar>
    </div>
    
    <div class="col-8 text-left">
      <div class="text-h6 q-pa-lg" style="margin: auto;">
        {{ userData?.name }}
        <q-btn v-if="isSameUser" @click="Edit" flat label="Edit"/>
        
        <!-- Follow/Unfollow buttons with loading state -->
        <q-btn 
          v-if="!isSameUser && !isUserFollowing"
          @click="FollowOrUnFollow" 
          flat 
          style="color: #FF0080" 
          :loading="followLoading"
          :disable="followLoading"
          label="Follow"
        />
        
        <q-btn 
          v-if="!isSameUser && isUserFollowing"
          @click="FollowOrUnFollow" 
          flat 
          class="primary" 
          :loading="followLoading"
          :disable="followLoading"
          label="Unfollow"
        />
      </div>
      
      <q-separator inset />
      
      <div class="text-subtitle1 q-pa-lg" style="margin: auto;">
        {{ userData?.bio }}
        <div>
          <i>{{ userPosts?.length || 0 }} Posts</i>
          <i>
            <i v-if="userData?.followers?.length > 0">
              {{ userData?.followers?.length }}
            </i>
            followers
          </i>
          <i>
            <i v-if="userData?.following?.length > 0">
              {{ userData?.following?.length }}
            </i>
            following
          </i>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'ShowProfile',
  props: ['userData', 'userPosts', 'isSameUser'],
  data() {
    return {
      isUserFollowing: false,
      followLoading: false
    }
  },
  
  watch: {
    // Watch userData to update follow status when it changes
    userData: {
      handler(newUserData) {
        if (newUserData && !this.isSameUser) {
          this.checkUserFollowingFromUserData();
        }
      },
      immediate: true
    },
    
    // Watch isSameUser to reset follow state
    isSameUser: {
      handler(newValue) {
        if (newValue) {
          this.isUserFollowing = false;
        }
      },
      immediate: true
    }
  },
  
  computed: {
    ...mapGetters(['GetUserData'])
  },
  
  methods: {
    ...mapActions(['FollowUser']),
    
    // Check follow status from userData prop instead of making API call
    checkUserFollowingFromUserData() {
      if (!this.userData || this.isSameUser) {
        this.isUserFollowing = false;
        return;
      }
      
      const loggedUserId = this.GetUserData()?.result?._id;
      
      if (!loggedUserId) {
        this.isUserFollowing = false;
        return;
      }
      
      // Check if logged user is in the followers array
      const followers = this.userData.followers || [];
      this.isUserFollowing = followers.some(followerId => 
        String(followerId) === String(loggedUserId)
      );
      
      console.log('Follow check - LoggedUser:', loggedUserId, 'IsFollowing:', this.isUserFollowing);
    },
    
    async FollowOrUnFollow() {
      if (this.isSameUser || this.followLoading || !this.userData?._id) {
        return;
      }
      
      try {
        this.followLoading = true;
        console.log("Following/Unfollowing user:", this.userData._id);
        
        const data = await this.FollowUser(this.userData._id);
        console.log("Follow response:", data);
        
        if (data) {
          // Update parent component with new user data
          if (data.FirstUser) {
            this.$emit('update-user', {
              data: data.FirstUser
            });
          }
          
          // Update local follow state immediately
          this.isUserFollowing = !this.isUserFollowing;
          
          // Show success notification
          this.$q.notify({
            type: 'positive',
            message: this.isUserFollowing ? 'Now following' : 'Unfollowed',
            timeout: 1500
          });
        }
        
      } catch (error) {
        console.error('Follow/Unfollow error:', error);
        this.$q.notify({
          type: 'negative',
          message: 'Failed to update follow status',
          timeout: 2000
        });
      } finally {
        this.followLoading = false;
      }
    },
    
    Edit() {
      this.$emit('EditProfile');
    }
  },
  
  mounted() {
    console.log('ShowProfile mounted - UserData:', this.userData, 'IsSameUser:', this.isSameUser);
    
    // Only check follow status if we have userData and it's not the same user
    if (this.userData && !this.isSameUser) {
      this.checkUserFollowingFromUserData();
    }
  }
}
</script>



<!-- <template>
 <div class="row col-12  ">
   <div class="col-4 text-center resize-observer-fix">
    <q-avatar size="150px">
     <img v-if="userData?.imageUrl" :src="userData?.imageUrl" loading="lazy">
     <img v-else src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png" loading="lazy">
   </q-avatar>
   </div>
   <div class="col-8 text-left">
    <div class="text-h6 q-pa-lg" style="margin: auto;">
        {{ userData?.name }}
        <q-btn v-if="isSameUser" @click="Edit" flat label="Edit"/>

         <q-btn v-if="!isSameUser && !isUserFollowing"
            @click="FollowOrUnFollow" flat style="color: #FF0080" label="Follow"/>

         <q-btn v-if="!isSameUser && isUserFollowing"
            @click="FollowOrUnFollow" flat class="primary" label="UN Follow"/>
    </div>
    <q-separator inset />
    <div class="text-subtitle1 q-pa-lg" style="margin: auto;">
        {{ userData.bio }}
        <div>
            <i>{{ userPosts.length }} Posts</i>
            <i>
                <i v-if="userData?.followers?.length > 0">
                    {{ userData?.followers?.length  }}</i>
                    followers
            </i>
            <i>
                <i v-if="userData?.following?.length > 0">
                    {{ userData?.following?.length  }}</i>
                    following
            </i>
        </div>
    </div>
   </div>
 </div>


</template>

<script>
import { mapActions } from 'vuex';
 export default {
    props:['userData','userPosts', 'isSameUser'],
    data(){
        return {isUserFollowing:false}
    }, 
    methods:{
        ...mapActions(['FollowUser', 'GetUserByID']),
        async checkUserFollowing(){
            const logeuid = JSON.parse(localStorage.getItem('profile'))?.result?._id
            // const id = this.userData?._id 

            const { user } = await this.GetUserByID(this.$route.params.id)

            if(user && user?.followers.find((id) => id == logeuid)){
                this.isUserFollowing = true 
            } else {
                this.isUserFollowing = false 
            }
            console.log("isUserFollo", this.isUserFollowing, "isSameUser", this.isSameUser)

        },
        async FollowOrUnFollow(){
            console.log("follow or un Follow user")
            let data = await this.FollowUser(this.userData._id)
            console.log("data show profile follow", data)
            if(data && data.FirstUser){

                this.$emit('update-user', {
                    data: data?.FirstUser
                })
                // change buttom
                this.checkUserFollowing()
            }
        },
        Edit(){
            this.$emit('EditProfile')
        },
    },
    mounted(){
        this.checkUserFollowing()
    }
 }
</script> -->