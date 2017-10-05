import { Component } from '@angular/core';
import { Menuitem } from '../../models/menuitem-models';
@Component({
    selector: 'nav-menu',
    templateUrl: './navmenu.component.html',
    styleUrls: ['./navmenu.component.css']
})
export class NavMenuComponent {
    public extraMenuItems: Array<Menuitem> = [{
        name:'styleguide',
        link:'/styleguide'
    }];
    updateMenuItems(value) {
        this.extraMenuItems = value;
      }
}
