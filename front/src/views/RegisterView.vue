<script setup lang="ts">
import { ref } from "vue";
// @ts-ignore
import { useStore } from 'vuex'
import { onMounted, computed } from 'vue'


import { REGISTER_ERROR_CLEAR, ACTION_REGISTER } from '@/store/auth'
import AuthenticationExplain from '@/components/AuthenticationExplain.vue';
import Form from '@/components/Form.vue'
import Error from '@/containers/Error.vue'
import Input from '@/containers/Input.vue'
import SubmitButton from '@/containers/SubmitButton.vue'

const store = useStore()

const user = ref("");
const group = ref("");
const displayName = ref("");

const userError = ref("")
const groupError = ref("")
const displayNameError = ref("")


const formError = computed(() => store.state.auth.registerError)


function onSubmit(e: Event) {
    e.preventDefault()
    let hasError = false
    if (user.value.length == 0) {
        hasError = true
        userError.value = "User name can't be empty"
    }

    if (group.value.length == 0) {
        hasError = true
        groupError.value = "Group name can't be empty"
    }

    if (displayName.value.length == 0) {
        hasError = true
        displayNameError.value = "User display name name can't be empty"
    }

    if (hasError) {
        return
    }

    store.dispatch(ACTION_REGISTER, {
        id: user.value,
        group_id: group.value,
        name: displayName.value,
    })
}

onMounted(() => {
    store.commit(REGISTER_ERROR_CLEAR)
})
</script>

<template>
    <div>
        <AuthenticationExplain />
    </div>

    <Form @submit="onSubmit">
        <Input placeholder="User" v-model="user" :error="userError"
            @focus="() => { userError = ''; store.commit(REGISTER_ERROR_CLEAR) }" />
        <Input placeholder="User display Name" v-model="displayName" :error="displayNameError"
            @focus="() => displayNameError = ''" />
        <Input placeholder="Group" v-model="group" :error="groupError" @focus="() => groupError = ''" />

        <Error :error="formError" />

        <SubmitButton name="Register" />
    </Form>
</template>

