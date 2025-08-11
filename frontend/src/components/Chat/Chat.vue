<template>
  <q-page class="constrain q-pa-md resize-observer-fix">
    <div class="row q-col-gutter-lg">
      <div class="col-12 chat-container">
        
        <!-- User list - keep it simple -->
        <div class="user-list">
          <div class="q-pa-md">
            <q-toolbar class="bg-primary text-white shadow-1">
              <q-toolbar-title>Following & Followers</q-toolbar-title>
            </q-toolbar>
            
            <q-list bordered>
              <q-item 
                @click="selectUser(contact)" 
                v-for="contact in contacts" 
                :key="contact._id" 
                class="q-my-sm" 
                clickable 
                v-ripple
                :class="{ 'bg-blue-1': selectedUser && selectedUser._id === contact._id }"
              >
                <q-item-section avatar>
                  <q-avatar v-if="!contact.imageUrl" color="primary" text-color="white">
                    {{ contact.name[0] }}
                  </q-avatar>
                  <q-avatar v-else>
                    <img :src="contact?.imageUrl">
                  </q-avatar>
                </q-item-section>
                
                <q-item-section>
                  <q-item-label>{{ contact.name }}</q-item-label>
                </q-item-section>

                <q-item-section side v-if="contact.isOnline">
                  <q-badge color="positive" rounded/>
                </q-item-section>

                <q-item-section 
                  side 
                  v-if="contact.unReadedmessage && contact.unReadedmessage > 0"
                >
                  <q-badge 
                    color="negative" 
                    rounded 
                    :label="contact?.unReadedmessage" 
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </div>

        <!-- Chat box with user header -->
        <div class="chat-messages" v-if="selectedUser != null" style="background: white;">
          
          <!-- Chat Header with selected user info -->
          <div class="chat-header">
            <q-toolbar class="bg-grey-1 text-primary shadow-1">
              <q-avatar class="q-mr-md">
                <img v-if="selectedUser.imageUrl" :src="selectedUser.imageUrl">
                <span v-else class="bg-primary text-white" style="width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center;">
                  {{ selectedUser.name[0] }}
                </span>
              </q-avatar>
              <q-toolbar-title>
                <div>
                  <div class="text-weight-medium">{{ selectedUser.name }}</div>
                  <div class="text-caption" v-if="selectedUser.isOnline">
                    <span class="online-indicator q-mr-xs"></span>
                    Online
                  </div>
                  <div class="text-caption text-grey" v-else>
                    <span class="offline-indicator q-mr-xs"></span>
                    Offline
                  </div>
                </div>
              </q-toolbar-title>
            </q-toolbar>
          </div>

          <!-- Messages Container -->
          <div 
            class="q-pa-md row justify-center messages-scroll-container" 
            style="overflow-y: auto; max-height: 400px;" 
            ref="messageContainer" 
            @scroll="handleScroll"
          >
            <!-- Loading indicator for fetching older messages -->
            <div v-if="loadingOlderMessages" class="full-width text-center q-py-md">
              <q-spinner color="primary" size="md" />
              <div class="text-caption q-mt-xs">Loading older messages...</div>
            </div>
            
            <div v-for="msg in messageBetweenUsers" :key="msg._id" style="width: 100%;">
              <q-chat-message 
                :name="msg.sender === MainUserData._id ? MainUserData.name : selectedUser.name" 
                :avatar="msg.sender === MainUserData._id ? MainUserData.imageUrl : selectedUser.imageUrl" 
                :text="[msg.content]" 
                :sent="msg.sender === MainUserData._id ? true : false" 
              />
            </div>
          </div>

          <q-separator spaced />
          
          <!-- Message Input -->
          <div class="q-pa-md">
            <q-input 
              outlined 
              v-model="messaageToSend.text" 
              @keyup.enter="Sendmessage" 
              label="write message.."
              :disable="!selectedUser"
            >
              <template v-slot:append>
                <q-btn 
                  v-if="messaageToSend.text != ''" 
                  @click="Sendmessage" 
                  flat 
                  round 
                  color="primary" 
                  icon="eva-arrow-right" 
                />
              </template>
            </q-input>
          </div>
        </div>

        <!-- Empty state when no user selected -->
        <div v-else class="chat-messages empty-state">
          <div class="text-center q-pa-xl">
            <q-icon name="eva-message-circle-outline" size="4rem" color="grey-4" />
            <div class="text-h6 text-grey-6 q-mt-md">Select a contact to start chatting</div>
            <div class="text-caption text-grey-5">Choose someone from your followers or following list</div>
          </div>
        </div>

      </div>
    </div>
  </q-page>
</template>

<script>
import { mapGetters, mapActions, mapState } from 'vuex';
export default {
    name:'ChatComponent',
    data(){
        return {
            messaageToSend: {text: ''},
            contacts:[],
            messageBetweenUsers:[],
            messagelistnum:0,
            selectedUser: null,
            MainUserData:{},
            uniqueOnlineUsers:[],
            loadingOlderMessages: false,
            hasMoreMessages: true,
            initialScrollDone: false,
        };
    },
    computed:{
        ...mapGetters(['GetUserFollowersFollowing','GetUserData']),
        ...mapState(["RealTimeChat"])
    },
    watch: {
        "RealTimeChat.onlineFriends": function (online) {
            const onlineFriendsArray = Object.values(online);
            console.log('Online friend changed new val', onlineFriendsArray)
            this.uniqueOnlineUsers = Array.from(new Set(onlineFriendsArray));
            this.updateOnlineList();
        },
        "RealTimeChat.privateMessages": function(message){
            if(this.contacts.length > 0){
                this.contacts.forEach((contact)=> {
                    if(contact._id == message.sender) {
                        contact.unReadedmessage++;
                    }
                })
                if(this.selectedUser && this.selectedUser?._id == message.sender){
                    this.messageBetweenUsers.push(message);
                    setTimeout(() => {
                        this.scrollDownFunction();
                    }, 100);
                }
            }
        }
    },
    async mounted(){
        this.MainUserData = this.GetUserData().result;
        this.GetUsList();
        this.uniqueOnlineUsers = Array.from(new Set(Object.values(this.RealTimeChat.onlineFriends)));
        this.updateOnlineList();
    },
    methods:{
        ...mapActions([
            'GetUnreadedMessageNum',
            'GetChatMsgsBetweenTwoUsers',
            'SendMessage',
            'MarkMsgsAsReaded',
            'SendPrivateMessage'
        ]),
        updateOnlineList(){
            this.contacts.forEach((contact)=> {
                if(this.uniqueOnlineUsers.includes(contact._id)) {
                    contact.isOnline = true;
                } else {
                    contact.isOnline = false;
                }
            })
        },
        handleScroll(){
            const container = this.$refs.messageContainer;
            if (!container || this.loadingOlderMessages || !this.hasMoreMessages) return;
            
            // Check if scrolled to top (with small threshold)
            if (container.scrollTop <= 10) {
                console.log('Scrolled to top, loading older messages...');
                this.GetOldestMessgesBetweenUsers();
            }
        },

        async GetOldestMessgesBetweenUsers(){
            if (this.loadingOlderMessages || !this.hasMoreMessages) return;
            
            this.loadingOlderMessages = true;
            this.messagelistnum = this.messagelistnum + 1;
            
            const container = this.$refs.messageContainer;
            const oldScrollHeight = container ? container.scrollHeight : 0;
            
            try {
                var firstuid = this.MainUserData._id
                var seconduid = this.selectedUser._id
                var from = this.messagelistnum;
                var ndata = {from, firstuid, seconduid};

                var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
                
                if (msgs && msgs.length > 0) {
                    this.messageBetweenUsers.unshift(...msgs);
                    
                    // Maintain scroll position after adding messages
                    this.$nextTick(() => {
                        if (container) {
                            const newScrollHeight = container.scrollHeight;
                            container.scrollTop = newScrollHeight - oldScrollHeight;
                        }
                    });
                } else {
                    // No more messages available
                    this.hasMoreMessages = false;
                }
            } catch (error) {
                console.error('Error loading older messages:', error);
            } finally {
                this.loadingOlderMessages = false;
            }
        },
        
        scrollDownFunction(){
            this.$nextTick(() => {
                const container = this.$refs.messageContainer;
                if (container) {
                    container.scrollTop = container.scrollHeight;
                }
            });
        },
        
        async CallMarkMsgAsReaded(user){
            var mainuid = this.MainUserData._id;
            var otheruid = user._id;
            var GetunReadedmessage = 0

            this.contacts.forEach(
                user => {
                    if(String(otheruid) == String(user._id)){
                        GetunReadedmessage = user.unReadedmessage
                    }
                }
            )

            var data = {mainuid, otheruid, GetunReadedmessage}
            var {isMarked} = await this.MarkMsgsAsReaded(data);

            if(isMarked){
                this.contacts.forEach(user => {
                    if(String(otheruid)== String(user._id)){
                        user.unReadedmessage = 0;
                    }
                })
            }
        },
        
        async GetUnreadedMsgList(){
            var {messages} = await this.GetUnreadedMessageNum(this.MainUserData._id);
            this.contacts.forEach(user => {
                messages.forEach(msg => {
                    if(String(msg.otherUserid) == String(user._id)){
                        user.unReadedmessage = Number(msg.numOfUnreadedMessages);
                    }
                })
            })
        },
        
        async GetUsList(){
            this.contacts = [];
            var glist = await this.GetUserFollowersFollowing;
            this.contacts = glist;
            if(this.contacts){
                this.GetUnreadedMsgList();
            }
            this.updateOnlineList();
        },
        
        async selectUser(user){
            // Reset states for new conversation
            this.selectedUser = null;
            this.messageBetweenUsers = [];
            this.messagelistnum = 0;
            this.hasMoreMessages = true;
            this.initialScrollDone = false;
            
            this.selectedUser = user;
            
            var firstuid = this.MainUserData._id;
            var seconduid = user._id;
            var from = 0;
            var ndata = {from, firstuid, seconduid};
            
            try {
                var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
                this.messageBetweenUsers.push(...msgs);
                
                // Scroll to bottom after initial load
                setTimeout(() => {
                    this.scrollDownFunction();
                    this.CallMarkMsgAsReaded(user);
                    this.initialScrollDone = true;
                }, 200);
                
            } catch (error) {
                console.error('Error loading initial messages:', error);
            }
        },
        
        Sendmessage(){
            if (!this.messaageToSend.text.trim() || !this.selectedUser) return;
            
            var content = this.messaageToSend.text;
            var sender = this.MainUserData._id;
            var recever = this.selectedUser._id;

            var sdata = {content, sender, recever};
            
            if(!this.uniqueOnlineUsers.includes(recever)) {
                var sucess = this.SendMessage(sdata);
                if (sucess){
                    this.messageBetweenUsers.push(sdata);
                    setTimeout(() => {
                        this.scrollDownFunction();
                    }, 100);
                }
            } else {
                this.SendPrivateMessage(sdata).then(()=> {
                    this.messageBetweenUsers.push(sdata);
                    setTimeout(() => {
                        this.scrollDownFunction();
                    }, 100);
                });
            }

            this.messaageToSend.text = '';
        }
    }
}
</script>

<style scoped>
.chat-container {
    display: flex;
}

.user-list {
    min-width: 300px;
    max-width: 350px;
    border-right: 1px solid #e0e0e0;
}

.chat-messages {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 500px;
}

.chat-header {
    border-bottom: 1px solid #e0e0e0;
}

.messages-scroll-container {
    flex: 1;
    background: #f5f5f5;
}

.empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fafafa;
}

/* Highlight selected user */
.bg-blue-1 {
    background-color: #e3f2fd !important;
}

/* Online/Offline indicators */
.online-indicator {
    display: inline-block;
    width: 8px;
    height: 8px;
    background-color: #4caf50;
    border-radius: 50%;
}

.offline-indicator {
    display: inline-block;
    width: 8px;
    height: 8px;
    background-color: #9e9e9e;
    border-radius: 50%;
}
.messages-scroll-container::-webkit-scrollbar {
    width: 6px;
}

.messages-scroll-container::-webkit-scrollbar-track {
    background: #f1f1f1;
}

.messages-scroll-container::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
}

.messages-scroll-container::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
}
</style>








<!-- <template>
  <q-page class="constrain q-pa-md resize-observer-fix">
    <div class="row q-col-gutter-lg">
      <div class="col-12 chat-container">
        
        <div class="user-list">
          <div class="q-pa-md">
            <q-toolbar class="bg-primary text-white shadow-1">
              <q-toolbar-title>Following & Followers</q-toolbar-title>
            </q-toolbar>
            
            <q-list bordered>
              <q-item 
                @click="selectUser(contact)" 
                v-for="contact in contacts" 
                :key="contact._id" 
                class="q-my-sm" 
                clickable 
                v-ripple
              >
                <q-item-section avatar>
                  <q-avatar v-if="!contact.imageUrl" color="primary" text-color="white">
                    {{ contact.name[0] }}
                  </q-avatar>
                  <q-avatar v-else>
                    <img :src="contact?.imageUrl">
                  </q-avatar>
                </q-item-section>
                
                <q-item-section>
                  <q-item-label>{{ contact.name }}</q-item-label>
                </q-item-section>

                <q-item-section side v-if="contact.isOnline">
                  <q-badge color="positive" rounded/>
                </q-item-section>

                <q-item-section 
                  side 
                  v-if="contact.unReadedmessage && contact.unReadedmessage > 0"
                >
                  <q-badge 
                    color="negative" 
                    rounded 
                    :label="contact?.unReadedmessage" 
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </div>

        <div class="chat-messages" v-if="selectedUser != null" style="background: white;">
          <div 
            class="q-pa-md row justify-center" 
            style="overflow-y: auto; max-height: 400px;" 
            ref="messageContainer" 
            @scroll="handleScroll"
          >
            <div v-for="msg in messageBetweenUsers" :key="msg._id" style="width: 100%;">
              <q-chat-message 
                :name="msg.sender === MainUserData._id ? MainUserData.name : selectedUser.name" 
                :avatar="msg.sender === MainUserData._id ? MainUserData.imageUrl : selectedUser.imageUrl" 
                :text="[msg.content]" 
                :sent="msg.sender === MainUserData._id ? true : false" 
              />
            </div>
          </div>

          <q-separator spaced />
          
          <q-input 
            outlined 
            v-model="messaageToSend.text" 
            @keyup.enter="Sendmessage" 
            label="write message.."
          >
            <template v-slot:append>
              <q-btn 
                v-if="messaageToSend.text != ''" 
                @click="Sendmessage" 
                flat 
                round 
                color="primary" 
                icon="eva-arrow-right" 
              />
            </template>
          </q-input>
        </div>

      </div>
    </div>
  </q-page>
</template>
<script>
import { mapGetters, mapActions, mapState } from 'vuex';
export default {
    name:'ChatComponent',
    data(){
        return {
            messaageToSend: {text: ''},
            contacts:[],
            messageBetweenUsers:[],
            messagelistnum:0,
            selectedUser: null,
            MainUserData:{},
            uniqueOnlineUsers:[],
        };
    },
    computed:{
        ...mapGetters(['GetUserFollowersFollowing','GetUserData']),
        ...mapState(["RealTimeChat"])
    },
    watch: {
        "RealTimeChat.onlineFriends": function (online) {
            const onlineFriendsArray = Object.values(online);
            console.log('Online friend changed new val', onlineFriendsArray)
                this.uniqueOnlineUsers = Array.from(new Set(onlineFriendsArray));
                this.updateOnlineList();
        },
        "RealTimeChat.privateMessages": function(message){
            if(this.contacts.length > 0){
                this.contacts.forEach((contact)=> {
                    if(contact._id == message.sender) {
                        contact.unReadedmessage++;
                    }
                })
            if(this.selectedUser && this.selectedUser?._id == message.sender){
                this.messageBetweenUsers.push(message);
                setTimeout(() => {
                    this.scrollDownFunction();
                }, 100);
            }

            }
        }
    },
    async mounted(){
        this.MainUserData = this.GetUserData().result;
        this.GetUsList();

        this.uniqueOnlineUsers = Array.from(new Set(Object.values(this.RealTimeChat.onlineFriends)));
        this.updateOnlineList();

    },
    methods:{
        ...mapActions([
            'GetUnreadedMessageNum',
            'GetChatMsgsBetweenTwoUsers',
            'SendMessage',
            'MarkMsgsAsReaded',
            'SendPrivateMessage'
        ]),
        updateOnlineList(){
            this.contacts.forEach((contact)=> {
                if(this.uniqueOnlineUsers.includes(contact._id)) {
                    contact.isOnline = true;
                } else {
                    contact.isOnline = false;
                }
            })
        },
        handleScroll(){
            const container = this.$refs.messageContainer;
            if (container.scrollTop === 0){
                // scorelled to the top
                this.GetOldestMessgesBetweenUsers();
            }},

        async GetOldestMessgesBetweenUsers(){
            this.messagelistnum = this.messagelistnum +1;
            var firstuid = this.MainUserData._id
            var seconduid = this.selectedUser._id
            var from = this.messagelistnum;
            var ndata = {from, firstuid, seconduid};

            var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
            this.messageBetweenUsers.unshift(...msgs);

        },
        scrollDownFunction(){
            const container = this.$refs.messageContainer;
            container.scrollTop = container.scrollHeight;
        },
        async CallMarkMsgAsReaded(user){
            var mainuid = this.MainUserData._id;
            var otheruid = user._id;
            var GetunReadedmessage = 0

            this.contacts.forEach(
                user => {
                    if(String(otheruid) == String(user._id)){
                        GetunReadedmessage = user.unReadedmessage
                    }
                }
            )

            var data = {mainuid, otheruid, GetunReadedmessage}
            var {isMarked} = await this.MarkMsgsAsReaded(data);

            if(isMarked){
                this.contacts.forEach(user => {
                    if(String(otheruid)== String(user._id)){
                        user.unReadedmessage = 0;
                    }
                })
            }
        },
        async GetUnreadedMsgList(){
            var {messages} = await this.GetUnreadedMessageNum(this.MainUserData._id);
            this.contacts.forEach(user => {
                messages.forEach(msg => {
                    if(String(msg.otherUserid) == String(user._id)){
                        user.unReadedmessage = Number(msg.numOfUnreadedMessages);
                    }
                })
            })
        },
        async GetUsList(){
            this.contacts = [];
            var glist = await this.GetUserFollowersFollowing;
            this.contacts = glist;
            if(this.contacts){
                this.GetUnreadedMsgList();
            }
            this.updateOnlineList();

        },
        async selectUser(user){
            this.selectedUser = null;
            this.messageBetweenUsers = [];

            this.selectedUser = user;
            this.messagelistnum = 0;
            var firstuid = this.MainUserData._id;
            var seconduid = user._id;
            var from = 0;
            var ndata = {from, firstuid, seconduid};
            var {msgs} = await this.GetChatMsgsBetweenTwoUsers(ndata);
            this.messageBetweenUsers.push(...msgs);
            setTimeout(() => {
                this.scrollDownFunction();
                this.CallMarkMsgAsReaded(user)
            }, 100);

        },
        Sendmessage(){
            var content = this.messaageToSend.text;
            var sender = this.MainUserData._id;
            var recever = this.selectedUser._id;

            var sdata = {content, sender, recever};
            if(!this.uniqueOnlineUsers.includes(recever)) {
                
                var sucess = this.SendMessage(sdata);
                if (sucess){
                    this.messageBetweenUsers.push(sdata);
                    setTimeout(() => {
                        this.scrollDownFunction();
                    }, 100);
                }
            } else {
                this.SendPrivateMessage(sdata).then(()=> {
                    this.messageBetweenUsers.push(sdata);
                });
                setTimeout(() => {
                        this.scrollDownFunction();
                }, 100);
            }


            this.messaageToSend.text = '';
        }
 
    }
}

</script>

<style scoped>
.chat-container {
    display: flex;
}

.chat-messages {
    flex: 1;
    padding: 10px;
}


</style>



 -->
