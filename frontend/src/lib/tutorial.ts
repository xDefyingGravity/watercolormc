/**
 * tutorial.ts
 *
 * Description:
 *   frontend — tutorial.ts module.
 *
 * Created: 8/11/25
 * Author: Will Ballantine
 *
 * @packageDocumentation
 * @copyright 2025-present Will Ballantine
 */

import Shepherd, { type Tour } from 'shepherd.js'
import 'shepherd.js/dist/css/shepherd.css'

let tour: Tour

export function startTour() {
	tour = new Shepherd.Tour({
		defaultStepOptions: {
			scrollTo: true,
			cancelIcon: {
				enabled: true
			},
		},

		useModalOverlay: true
	});

	tour.addStep({
		id: 'step-1',
		text: 'Welcome to the server management interface! Here you can manage your Minecraft servers.',
		attachTo: {
			element: '',
			on: 'bottom'
		},
		buttons: [
			{
				text: 'Next',
				action: tour.next
			}
		]
	});

	tour.addStep({
		id: 'step-2',
		text: 'To create a new server, click the "+" button at the bottom right corner.',
		attachTo: {
			element: '#add-button',
			on: 'right'
		},
		buttons: [
			{
				text: 'Back',
				action: tour.back
			},
			{
				text: 'Next',
				action: tour.next
			}
		]
	});

tour.addStep({
		id: 'step-3',
		text: 'You can also create a server from the "Create" tab in the top navigation bar.',
		attachTo: {
			element: '#create-button',
			on: 'left'
		},
		buttons: [
			{
				text: 'Back',
				action: tour.back
			},
			{
				text: 'Finish',
				action: () => {
					tour.complete();
					localStorage.setItem('tutorial', 'completed');
				}
			}
		]
	});

	tour.start().catch(console.error);
}