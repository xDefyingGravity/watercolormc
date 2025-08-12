/**
 * validate.ts
 *
 * Description:
 *   frontend — validate.ts module.
 *
 * Created: 7/10/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */

import Ajv, { type ErrorObject } from 'ajv'
import schema from './schema.json'

export interface InvalidProperty {
	key: string
	message: string
}

const ajv = new Ajv({ allErrors: true, strict: false })
const validate = ajv.compile(schema)

function coerceValue(key: string, value: string): unknown {
	const def = schema.properties?.[key as keyof typeof schema.properties]
	if (!def) return value

	switch (def.type) {
		case 'boolean': {
			if (value === 'true') return true
			if (value === 'false') return false
			return '__invalid_boolean__'
		}
		case 'integer': {
			const parsed = parseInt(value, 10)
			if (value.trim() === '' || isNaN(parsed) || !Number.isInteger(parsed))
				return '__invalid_integer__'
			return parsed
		}
		case 'number': {
			const parsed = parseFloat(value)
			if (value.trim() === '' || isNaN(parsed)) return '__invalid_number__'
			return parsed
		}
		default:
			return value
	}
}

export function validateServerProperties(props: Record<string, string>): InvalidProperty[] {
	const coerced: Record<string, unknown> = {}
	for (const [key, value] of Object.entries(props)) {
		const v = coerceValue(key, value)
		if (v === '__invalid_boolean__') {
			return [{ key, message: 'Must be true or false' }]
		}
		if (v === '__invalid_integer__') {
			return [{ key, message: 'Must be a valid integer' }]
		}
		if (v === '__invalid_number__') {
			return [{ key, message: 'Must be a valid number' }]
		}
		coerced[key] = v
	}
	const valid = validate(coerced)
	if (valid) return []

	return (validate.errors ?? []).map((err: ErrorObject) => {
		const key = err.instancePath.startsWith('/')
			? err.instancePath.slice(1)
			: err.instancePath || 'root'
		const message = err.message ?? 'invalid value'
		return { key, message }
	})
}
