import { Component } from "@angular/core";
import {FormsModule} from "@angular/forms";
     
@Component({
    selector: "modbus-emulator",
    standalone: true,
    imports: [FormsModule],
    template: `<label>Введите ID сервера:</label>
                 <input [(ngModel)]="id" placeholder="id" type="number" min="0">
                 <p>Текущие данные для {{id}}...</p>`
})
export class AppComponent { 
    id = "";
}