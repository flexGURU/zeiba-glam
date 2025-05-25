import { Component } from '@angular/core';
import { Menubar } from 'primeng/menubar';

@Component({
  selector: 'app-menu',
  imports: [Menubar],
  templateUrl: './menu.component.html',
  styleUrl: './menu.component.css',
})
export class MenuComponent {
  items: any;
  baseLink: string = 'admin';

  ngOnInit() {
    this.items = [
      {
        label: 'Dashboard',
        icon: 'pi pi-home',
        routerLink: 'dashboard',
      },
      {
        label: 'Product list',
        icon: 'pi pi-tags',
        routerLink: 'catalog',
      },
    ];
  }
}
