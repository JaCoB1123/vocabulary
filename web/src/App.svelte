<script>
  import { Router } from "@sveltech/routify";
  import { routes } from "../.routify/routes";

  import { metatags } from '@sveltech/routify'
  metatags.title = 'Vocabulary'
  metatags.description = 'Simple vocabulary app'
</script>

<style  global>
  @import "../static/global.css";
</style>

<Router {routes} />