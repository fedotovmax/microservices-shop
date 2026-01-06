import Alpine from 'alpinejs';
import 'htmx.org';

import { A } from '../shared/keys';

function Hello() {
	console.log(`Variable from other file: ${A}`);
}

Alpine.start();

Hello();
