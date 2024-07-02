<template>
  <component :is="type" v-bind="linkProps(to,menuName)" >
    <slot />
  </component>
</template>

<script>
import { isExternal } from '@/utils/validate'

export default {
  props: {
    menuName: {
      type: String,
      required: true
    },
    to: {
      type: String,
      required: true
    }
  },
  computed: {
    isExternal() {
      return isExternal(this.to)
    },
    type() {
      if (this.isExternal) {
        return 'a'
      }
      return 'router-link'
    }
  },
  methods: {
    linkProps(to, menuName) {
      if (this.isExternal) {
        return {
          href: to,
          rel: 'noopener',
          target: '_blank'
        }
      }
      return {
        to: to
      }
    }
  }
}
</script>
