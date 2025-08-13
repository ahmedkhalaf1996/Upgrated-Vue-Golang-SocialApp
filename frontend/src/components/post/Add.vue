<template>
    <q-page-sticky 
        :position="$q.screen.lt.sm ? 'bottom-right' : 'bottom-left'"  
        v-show="GetUserData()?.result"
        :offset="$q.screen.lt.sm ? [18, 18] : [18, 60]"
    >
        <div class="q-pa-md q-gutter-sm">
            <q-btn 
                :label="$q.screen.gt.xs ? 'Create Post' : ''" 
                style="cursor: pointer;" 
                icon="eva-plus-circle-outline" 
                color="primary" 
                @click="persistent = true"
                :round="$q.screen.lt.sm"
                :size="$q.screen.lt.sm ? 'lg' : 'md'"
                class="touch-friendly"
            />

            <!-- Mobile-Optimized Dialog -->
            <q-dialog 
                v-model="persistent" 
                persistent 
                transition-show="slide-up" 
                transition-hide="slide-down"
                :maximized="$q.screen.lt.sm"
                :position="$q.screen.lt.sm ? 'bottom' : 'standard'"
            >
                <q-card :style="$q.screen.lt.sm ? 'min-height: 80vh; border-radius: 16px 16px 0 0;' : 'min-width: 350px;'">
                    <!-- Mobile Header -->
                    <q-card-section class="row items-center q-pb-none" v-if="$q.screen.lt.sm">
                        <div class="text-h6">Create Post</div>
                        <q-space />
                    </q-card-section>
                    
                    <!-- Desktop Header -->
                    <q-card-section v-else>
                        <div class="text-h6">Create Post</div>
                    </q-card-section>

                    <q-card-section class="q-pt-none">
                        <q-input 
                            dense 
                            v-model="post.title" 
                            autofocus 
                            placeholder="Post Title"
                            class="q-mb-md"
                        />
                        
                        <div class="q-mb-md">
                            <q-input
                                v-model="post.message"
                                placeholder="What's on your mind?"
                                type="textarea"
                                :rows="$q.screen.lt.sm ? 4 : 3"
                                autogrow
                            />
                        </div>
                        
                        <div class="q-mb-md">
                            <q-file 
                                v-model="file"
                                label="Pick Image"
                                filled 
                                accept="image/*"
                                :style="$q.screen.lt.sm ? 'width: 100%;' : 'max-width: 400px;'"
                            >
                                <template v-slot:prepend>
                                    <q-icon name="eva-camera-outline" />
                                </template>
                            </q-file>
                        </div>

                        <div class="q-gutter-sm row items-start" v-if="post.selectedFile">
                            <q-img 
                                :src="post.selectedFile"
                                spinner-color="red"
                                :style="$q.screen.lt.sm ? 'height: 200px; max-width: 100%;' : 'height: 140px; max-width: 150px;'"
                                class="rounded-borders"
                            />
                        </div>
                    </q-card-section>

                    <!-- Mobile Actions (Full Width) -->
                    <q-card-actions v-if="$q.screen.lt.sm" class="q-pa-md">
                        <q-btn 
                            flat 
                            label="Cancel" 
                            v-close-popup 
                            class="col"
                            size="md"
                        />
                        <q-btn 
                            unelevated
                            label="Create" 
                            color="primary"
                            v-close-popup 
                            @click="CreatePost"
                            class="col"
                            size="md"
                        />
                    </q-card-actions>
                    
                    <!-- Desktop Actions -->
                    <q-card-actions v-else align="right" class="text-primary">
                        <q-btn flat label="Cancel" v-close-popup />
                        <q-btn flat label="Create" v-close-popup @click="CreatePost"/>
                    </q-card-actions>
                </q-card>
            </q-dialog>
        </div>
    </q-page-sticky>  
</template>

<script>

import {mapActions, mapGetters} from 'vuex'

export default {
    name: 'AddComponent',
    data (){
      return {
        persistent: false,
        post: {title:'', message:'', name:'', selectedFile: null},
        file: null
      }        
    },
    watch:{
        file(){
            // Only convert if file exists and is valid
            this.ConvertToBase64();
        }
    },
    computed: {
        ...mapGetters(['GetUserData'])
    },
    methods: {
        ...mapActions(['createPost']),
        
        async CreatePost(){
            var name = JSON.parse(localStorage.getItem('profile'))?.result?.name;
            this.post.name = name;
            
            // Enhanced validation
            var isValidate = true;
            
            // Check title
            if (!this.post.title || this.post.title.trim() === '') {
                this.$q.notify({
                    icon: 'eva-alert-circle-outline',
                    type: 'negative',
                    message: 'Title is required'
                });
                isValidate = false;
            }
            
            // Check message with minimum length requirement
            if (!this.post.message || this.post.message.trim() === '') {
                this.$q.notify({
                    icon: 'eva-alert-circle-outline',
                    type: 'negative',
                    message: 'Message is required'
                });
                isValidate = false;
            } else if (this.post.message.trim().length < 5) {
                this.$q.notify({
                    icon: 'eva-alert-circle-outline',
                    type: 'negative',
                    message: 'Message must be at least 5 characters long'
                });
                isValidate = false;
            }
            
            // Check name
            if (!name) {
                this.$q.notify({
                    icon: 'eva-alert-circle-outline',
                    type: 'negative',
                    message: 'User name is required'
                });
                isValidate = false;
            }
            
            // Create post if validation passes
            if (isValidate) {
                try {
                    // Prepare post data - trim whitespace
                    const postData = {
                        title: this.post.title.trim(),
                        message: this.post.message.trim(),
                        selectedFile: this.post.selectedFile || ''  // Ensure it's never null
                    };
                    
                    const data = await this.createPost(postData);
                    console.log('data', data);

                    if (data) {
                        // Reset form - do this carefully to avoid watcher issues
                        this.resetForm();
                        
                        this.$emit('Created');
                        
                        this.$q.notify({
                            icon: 'eva-checkmark-circle-outline',
                            type: 'positive',
                            message: 'Post Created Successfully'
                        });
                    }
                } catch (error) {
                    console.error('Error creating post:', error);
                    
                    // Handle specific validation errors from backend
                    if (error.response && error.response.status === 400) {
                        const errorData = error.response.data;
                        if (Array.isArray(errorData)) {
                            // Handle validation errors
                            errorData.forEach(err => {
                                let message = `${err.Field} validation failed`;
                                if (err.Tag === 'min') {
                                    message = `${err.Field} must be at least 5 characters long`;
                                } else if (err.Tag === 'required') {
                                    message = `${err.Field} is required`;
                                }
                                
                                this.$q.notify({
                                    icon: 'eva-alert-circle-outline',
                                    type: 'negative',
                                    message: message
                                });
                            });
                        } else {
                            this.$q.notify({
                                icon: 'eva-alert-circle-outline',
                                type: 'negative',
                                message: 'Failed to create post. Please check your input.'
                            });
                        }
                    } else {
                        this.$q.notify({
                            icon: 'eva-alert-circle-outline',
                            type: 'negative',
                            message: 'An error occurred while creating the post'
                        });
                    }
                }
            }
        },
        
        resetForm() {
            // Reset form data carefully to avoid triggering watchers with invalid values
            this.post = {
                title: '', 
                message: '', 
                name: '', 
                selectedFile: null
            };
            // Reset file separately to avoid watcher issues
            this.$nextTick(() => {
                this.file = null;
            });
        },
        
        ConvertToBase64(){
            // Check if file exists and is a valid File/Blob object
            if (!this.file || !(this.file instanceof File || this.file instanceof Blob)) {
                // If file is null or not a valid file, reset selectedFile
                this.post.selectedFile = null;
                return;
            }

            const reader = new FileReader();
            
            reader.onload = () => {
                this.post.selectedFile = reader.result;
            };
            
            reader.onerror = (error) => {
                console.error('Error reading file:', error);
                this.$q.notify({
                    icon: 'eva-alert-circle-outline',
                    type: 'negative',
                    message: 'Error reading file'
                });
                this.post.selectedFile = null;
            };
            
            reader.readAsDataURL(this.file);
        }
    }
}
</script>


