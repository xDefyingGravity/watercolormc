import log from 'loglevel'

const DEBUG_COLOR = 'color: #00f; font-weight: bold;'
const INFO_COLOR = 'color: #0a0; font-weight: bold;'
const WARN_COLOR = 'color: #fa0; font-weight: bold;'
const ERROR_COLOR = 'color: #f00; font-weight: bold;'

export function init() {
	log.setLevel('debug')
}

export function debug(...message: any[]) {
	log.debug('%c[debug]', DEBUG_COLOR, ...message)
}

export function info(...message: any[]) {
	log.info('%c[info]', INFO_COLOR, ...message)
}

export function warn(...message: any[]) {
	log.warn('%c[warn]', WARN_COLOR, ...message)
}

export function error(...message: any[]) {
	log.error('%c[error]', ERROR_COLOR, ...message)
}
