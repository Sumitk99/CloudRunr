import { Routes } from '@angular/router';
import { LandingPageComponent } from './components/landing-page/landing-page.component';
import { AuthComponent } from './components/auth/auth.component';
import { NewProjectComponent } from './components/projects/new-project/new-project.component';
import { DeploymentPageComponent } from './components/projects/deployment/deployment-page.component';
import { ProjectsListComponent } from './components/projects/projects-list/projects-list.component';
import { ProjectDetailComponent } from './components/projects/project-detail/project-detail.component';

export const routes: Routes = [
  { path: '', component: LandingPageComponent },
  { path: 'auth', component: AuthComponent },
  { path: 'projects', component: ProjectsListComponent },
  { path: 'projects/new', component: NewProjectComponent },
  { path: 'projects/:projectId', component: ProjectDetailComponent },
  { path: 'projects/:projectId/:deploymentId', component: DeploymentPageComponent },
  { path: '**', redirectTo: '/', pathMatch: 'full' }
];
