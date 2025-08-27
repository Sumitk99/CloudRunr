import { Routes } from '@angular/router';
import { LandingPageComponent } from './components/landing-page/landing-page.component';

export const routes: Routes = [
  { path: '', component: LandingPageComponent },
  { path: 'login', redirectTo: '/', pathMatch: 'full' }, // Will implement login page later
  { path: '**', redirectTo: '/', pathMatch: 'full' }
];
