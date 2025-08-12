<script lang="ts">
	import { marked } from 'marked'
	import DOMPurify from 'dompurify'

	export let markdown = ''

	$: raw = marked.parse(markdown)
	$: safe = DOMPurify.sanitize(raw as string)
</script>

<div class="markdown-body max-w-none">
	{@html safe}
</div>
