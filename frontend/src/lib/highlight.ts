function escapeHtml(str: string) {
	return str
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
		.replace(/"/g, '&quot;')
		.replace(/'/g, '&#39;')
}

export function highlight(log: string) {
	if (!log) return log

	log = escapeHtml(log)

	log = log.replace(/^\[(.*?)\]/, '<span class="syntax__timestamp selectable">[$1]</span>')

	log = log.replace(
		/\[([^\]]+?)\/(ERROR|WARN|INFO|DEBUG|TRACE)\]/g,
		'<span class="syntax__thread selectable">[$1/</span><span class="syntax_$2">$2</span><span class="log_thread">]</span>'
	)

	return log
}
