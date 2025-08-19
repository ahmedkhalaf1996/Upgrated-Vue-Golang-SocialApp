<template>
  <q-page class="constrain q-pa-md resize-observer-fix">
    <!-- Mobile Toggle Tabs (only visible on mobile) -->
    <div class="mobile-tabs q-mb-md gt-xs-hide">
      <q-tabs
        v-model="activeTab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="justify"
        narrow-indicator
      >
        <q-tab name="signin" label="Sign In" />
        <q-tab name="signup" label="Sign Up" />
      </q-tabs>
    </div>

    <!-- Desktop Layout (hidden on mobile) -->
    <div class="row q-col-gutter-lg lt-sm-hide">
      <div class="col-5">
        <q-card class="my-card" style="width: 100%; padding: 10px;">
          <h1 class="text-h6 text-center">Sign In</h1>
          <q-card-section>
            <form @submit.prevent.stop="Login" class="q-gutter-md">
              <q-input
                filled 
                v-model="Sin_data.email"
                label="Your Email *"
                hint="Email"
                type="email"
                lazy-rules
              />
              <q-input
                filled 
                v-model="Sin_data.password"
                label="Your Password *"
                hint="Password"
                type="password"
                lazy-rules
              />  
              <div>
                <q-btn 
                  label="Sign In" 
                  type="submit" 
                  color="primary" 
                  class="full-width"
                  :loading="isLoading"
                />
              </div>                      
            </form>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-7">
        <q-card class="my-card" style="width: 100%; padding: 10px;">
          <h1 class="text-h6 text-center">Sign Up | Create New Account</h1>
          <q-card-section>
            <form @submit.prevent.stop="Register" class="q-gutter-md">
              <q-input
                filled 
                v-model="Sup_data.firstName"
                label="Your First Name *"
                hint="First Name"
                lazy-rules
              />
              <q-input
                filled 
                v-model="Sup_data.lastName"
                label="Your Last Name *"
                hint="Last Name"
                lazy-rules
              />
              <q-input
                filled 
                v-model="Sup_data.email"
                label="Your Email *"
                hint="Email"
                type="email"
                lazy-rules
              />
              <q-input
                filled 
                v-model="Sup_data.password"
                type="password"
                label="Your Password *"
                hint="Password"
                lazy-rules
              />  
              <div>
                <q-btn 
                  label="Create New Account" 
                  type="submit" 
                  color="positive" 
                  class="full-width"
                  :loading="isLoading"
                />
              </div>                      
            </form>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Mobile Layout (only visible on mobile) -->
    <div class="mobile-layout gt-xs-hide">
      <!-- Sign In Tab -->
      <q-card v-show="activeTab === 'signin'" class="mobile-card">
        <q-card-section class="text-center q-pb-none">
          <div class="text-h5 q-mb-xs">Welcome Back</div>
          <div class="text-grey-6">Sign in to your account</div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <form @submit.prevent.stop="Login" class="q-gutter-md q-mt-md">
            <q-input
              filled 
              v-model="Sin_data.email"
              label="Email Address"
              type="email"
              lazy-rules
              :rules="[val => !!val || 'Email is required']"
              class="mobile-input"
            />
            <q-input
              filled 
              v-model="Sin_data.password"
              label="Password"
              type="password"
              lazy-rules
              :rules="[val => !!val || 'Password is required']"
              class="mobile-input"
            />  
            <div class="q-mt-lg">
              <q-btn 
                label="Sign In" 
                type="submit" 
                color="primary" 
                class="full-width mobile-btn"
                :loading="isLoading"
                size="md"
                rounded
              />
            </div>                      
          </form>
        </q-card-section>
      </q-card>

      <!-- Sign Up Tab -->
      <q-card v-show="activeTab === 'signup'" class="mobile-card">
        <q-card-section class="text-center q-pb-none">
          <div class="text-h5 q-mb-xs">Join Us</div>
          <div class="text-grey-6">Create your new account</div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <form @submit.prevent.stop="Register" class="q-gutter-md q-mt-md">
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input
                  filled 
                  v-model="Sup_data.firstName"
                  label="First Name"
                  lazy-rules
                  :rules="[val => !!val || 'First name is required']"
                  class="mobile-input"
                />
              </div>
              <div class="col-6">
                <q-input
                  filled 
                  v-model="Sup_data.lastName"
                  label="Last Name"
                  lazy-rules
                  :rules="[val => !!val || 'Last name is required']"
                  class="mobile-input"
                />
              </div>
            </div>
            <q-input
              filled 
              v-model="Sup_data.email"
              label="Email Address"
              type="email"
              lazy-rules
              :rules="[val => !!val || 'Email is required']"
              class="mobile-input"
            />
            <q-input
              filled 
              v-model="Sup_data.password"
              type="password"
              label="Password"
              lazy-rules
              :rules="[val => !!val || 'Password is required']"
              class="mobile-input"
            />  
            <div class="q-mt-lg">
              <q-btn 
                label="Create Account" 
                type="submit" 
                color="positive" 
                class="full-width mobile-btn"
                :loading="isLoading"
                size="md"
                rounded
              />
            </div>                      
          </form>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>
  
<script>
import {mapActions} from 'vuex'

export default {
  name: 'AuthView',
  data () {
    return {
        activeTab: 'signin', // For mobile tab switching
        isLoading: false,
        Sin_data:{
            email:'',
            password: '',
        },
        Sup_data:{
            email:'',
            password: '',
            firstName: '',
            lastName:'',
        }
    }
  },
  methods:{
    ...mapActions([
        'signin', 
        'signup',
        'connectToNotify',
        'createChatConnection'
    ]),
    async Login(){
        console.log("login in data", this.Sin_data)
        this.isLoading = true;
        var validate = true 
        
        if (this.Sin_data.email == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Email is required`,
                position: 'top'
           })    
        } else if (this.Sin_data.password == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Password is required`,
                position: 'top'
           })  
        }
        
        if(validate){
            try {
                var formdata = {email:this.Sin_data.email, password: this.Sin_data.password};
                const data  = await this.signin(formdata);

                if(data?.response?.data?.message || data?.response?.data){
                    this.$q.notify({
                        icon:'eva-alert-circle-outline',
                        type:'negative',
                        message:`Error: ${data.response.data.message ?  data?.response?.data.message :  data?.response?.data}`,
                        position: 'top'
                    })  
                } else {
                    this.$q.notify({
                        icon:'eva-checkmark-circle-outline',
                        type:'positive',
                        message:`Successfully signed in`,
                        position: 'top'
                    })  
                    this.connectToNotify()
                    this.createChatConnection()
                    this.$router.push('/')
                }
            } catch (error) {
                this.$q.notify({
                    icon:'eva-alert-circle-outline',
                    type:'negative',
                    message:`An error occurred during sign in`,
                    position: 'top'
                })
            }
        }
        this.isLoading = false;
    },
    async Register(){
        console.log("Register in data", this.Sup_data)
        this.isLoading = true;
        
        var isValidate = true;
        for (const key in this.Sup_data){
            const val = this.Sup_data[key];
            if(val === ''){
                this.$q.notify({
                    icon:'eva-alert-circle-outline',
                    type:'negative',
                    message:`${key.charAt(0).toUpperCase() + key.slice(1)} is required`,
                    position: 'top'
                });
                isValidate = false
                break;
            }
        }
        
        if(isValidate){
            try {
                const data = await this.signup(this.Sup_data)
                console.log("data on Register", data)
                
                if(data?.response?.data?.message){
                    this.$q.notify({
                        icon:'eva-alert-circle-outline',
                        type:'negative',
                        message:`Error: ${data.response.data.message}`,
                        position: 'top'
                    })  
                } else {
                    this.$q.notify({
                        icon:'eva-checkmark-circle-outline',
                        type:'positive',
                        message:`Successfully signed up`,
                        position: 'top'
                    })  
                    this.connectToNotify()
                    this.createChatConnection()
                    // Optionally switch to sign in tab or redirect
                    // this.activeTab = 'signin';
                    // this.$router.push('/')
                }
            } catch (error) {
                this.$q.notify({
                    icon:'eva-alert-circle-outline',
                    type:'negative',
                    message:`An error occurred during sign up`,
                    position: 'top'
                })
            }
        }
        this.isLoading = false;
    },
  }
}
</script>

<style lang="scss" scoped>
.mobile-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  max-width: 400px;
  margin: 0 auto;
}

.mobile-input {
  .q-field__control {
    border-radius: 12px;
  }
}

.mobile-btn {
  height: 48px;
  font-weight: 600;
  text-transform: none;
  letter-spacing: 0.5px;
}

.mobile-tabs {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 4px;
}

// Desktop improvements
.my-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

// Responsive adjustments
@media (max-width: 599px) {
  .constrain {
    padding: 16px 8px !important;
  }
}

// Hide elements based on screen size
.gt-xs-hide {
  @media (min-width: 600px) {
    display: none !important;
  }
}

.lt-sm-hide {
  @media (max-width: 599px) {
    display: none !important;
  }
}
</style>





<!-- <template>
 <q-page class="constrain q-pa-md resize-observer-fix">
    <div class="row q-col-gutter-lg">
        <div class="col-5">
            <q-card class="my-card" style="width: 100%; padding: 10px;">
                <h1 class="text-h6 text-center">Signin</h1>
                <q-card-section>
                    <form @submit.prevent.stop="Login" class="q-gutter-md">
                        <q-input
                          filled 
                           v-model="Sin_data.email"
                           label="Your Email *"
                           hint="Email"
                           lazy-rules
                        />
                        <q-input
                           filled 
                           v-model="Sin_data.password"
                           label="Your Password *"
                           hint="password"
                           type="password"
                           lazy-rules
                        />  
                        <div>
                            <q-btn label="sigin in" type="submit" color="primary" />
                        </div>                      
                    </form>
                </q-card-section>
            </q-card>
        </div>
        <div class="col-7">
            <q-card class="my-card" style="width: 100%; padding: 10px;">
                <h1 class="text-h6 text-center">Signup | Craete New Account</h1>
                <q-card-section>
                    <form @submit.prevent.stop="Register" class="q-gutter-md">
                        <q-input
                          filled 
                           v-model="Sup_data.firstName"
                           label="Your first Name *"
                           hint="firstName"
                           lazy-rules
                        />
                        <q-input
                          filled 
                           v-model="Sup_data.lastName"
                           label="Your lastName *"
                           hint="lastName"
                           lazy-rules
                        />
                        <q-input
                          filled 
                           v-model="Sup_data.email"
                           label="Your Email *"
                           hint="Email"
                           lazy-rules
                        />
                        <q-input
                           filled 
                           v-model="Sup_data.password"
                           type="password"
                           label="Your Password *"
                           hint="password"
                           lazy-rules
                        />  
                        <div>
                            <q-btn label="Create New Account" type="submit" color="positive" />
                        </div>                      
                    </form>
                </q-card-section>
            </q-card>
        </div>
    </div>
 </q-page>
</template>
  
<script>

import {mapActions} from 'vuex'

export default {
  name: 'AuthView',
  data () {
    return {
        Sin_data:{
            email:'',
            password: '',
        },
        Sup_data:{
            email:'',
            password: '',
            firstName: '',
            lastName:'',
        }
    }
  },
  methods:{
    ...mapActions([
        'signin', 
        'signup',
        'connectToNotify',
        'createChatConnection'
    ]),
    async Login(){
        console.log("login in data", this.Sin_data)
        var validate = true 
        if (this.Sin_data.email == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Email is Required`
           })    
        } else if (this.Sin_data.password == ''){
            validate = false 
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`password is Required`
           })  
        }
        // sucess and ready to next
        if(validate){
            var formdata = {email:this.Sin_data.email, password: this.Sin_data.password};
            const data  = await this.signin(formdata);
            // console.log("data", data, 'message', message)
            // console.log("response data on Login", data)
            // console.log('data response', data.response.data.message )

            if(data?.response?.data?.message || data?.response?.data){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Erorr ${data.response.data.message ?  data?.response?.data.message :  data?.response?.data}`
           })  
          } else {
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'positive',
                message:`Successfully Sigin in`
           })  
           this.connectToNotify()
           this.createChatConnection()
           this.$router.push('/')
          }
        }

    },
    async Register(){
        console.log("Register in data", this.Sup_data)
        // validation
        var isVaidate = true;
        for (const key in this.Sup_data){
            const val = this.Sup_data[key];
            if(val === ''){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`${key} is Required`
              });
              isVaidate = false
            }
        } // v end
        if(isVaidate){
            const data = await this.signup(this.Sup_data)

            console.log("data on Register", data)
            if(data?.response?.data?.message){
                this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'negative',
                message:`Erorr ${data.response.data.message}`
           })  
          } else {
            // meaning succusfully
            this.$q.notify({
                icon:'eva-alert-circle-outline',
                type:'positive',
                message:`Successfully Sigin up`
           })  
        this.connectToNotify()
        this.createChatConnection()

        //    this.$router.push('/')
          }
        }
    },
  }
}
</script> -->