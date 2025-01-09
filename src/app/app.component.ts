import { Component } from "@angular/core";
import {FormsModule} from "@angular/forms";
import { HttpClient, HttpClientModule } from "@angular/common/http";
import { ActualTime } from "./actual_time"
     
@Component({
    selector: "modbus-emulator",
    standalone: true,
    imports: [FormsModule, HttpClientModule],
    template: ` <h1>Modbus util: emulator</h1>
                <label>Enter server's ID:</label>
                <input [(ngModel)]="server_id" type="number" min="-1">
                <p>if server ID == -1 -> actions will done for all servers
                <p>Status: {{status}}</p>
                <button (click)="getActualTime()">Get actual time</button>
                @for(current_data of actual_times; track $index){
                    <li>{{current_data.id}} - {{current_data.actual_time}}</li>
                } @empty {
                    <li>There are no data</li>
                }
                <p>
                    <input type="datetime-local" step="1" [(ngModel)]="rewind_time">
                    <button (click)="rewindEmulation()">Set emulation time</button>
                </p>`
})

export class AppComponent{  
    actual_times: ActualTime[]=[];
    server_id: number = -1
    status: string = "Waiting"
    rewind_time: string
    constructor(private http: HttpClient) {}
    getActualTime(): void {
        let request = 'modbus-emulator/time/actual'
        if (this.server_id!==-1) {
            request = request + `?server_id=${this.server_id}`
        }
        if (this.actual_times.length !== 0 && this.server_id > this.actual_times.length) {
            console.error("Error: server ID must be in [0:%d]", (this.actual_times.length + 1))
            this.status = `Error - server ID must be in [0:${this.actual_times.length + 1}]`
            return
        }
        this.http.get(request).subscribe({next:(data:any) => this.actual_times=data, error: error => console.log(error)});
        this.status = "Success"
    }
    rewindEmulation(): void {
        this.rewind_time = this.rewind_time.replace("T", "%20")
        this.rewind_time = this.rewind_time.replaceAll(":", "%3A")
        console.log(this.rewind_time)
        let request = `modbus-emulator/time/rewind_emulation?timepoint=${this.rewind_time}`
        if (this.server_id!==-1) {
            request = request + `&server_id=${this.server_id}`
        }
        this.http.post(request, {}).subscribe({error: error => console.log(error)});
        this.status = "Success"
    }
}