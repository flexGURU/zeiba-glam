import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MenuComponent } from '../menu/menu.component';
import { AdminHeaderComponent } from '../admin-header/admin-header.component';
import { DashboardComponent } from '../../components/dashboard/dashboard.component';

@Component({
  selector: 'app-admin-layout',
  imports: [RouterOutlet, MenuComponent, AdminHeaderComponent],
  templateUrl: './admin-layout.component.html',
  styleUrl: './admin-layout.component.css',
})
export class AdminLayoutComponent {}
