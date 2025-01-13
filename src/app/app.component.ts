import { Component } from "@angular/core";
import {FormsModule} from "@angular/forms";
import { HttpClient, HttpClientModule, HttpHandler } from "@angular/common/http";
import { interval, Subscription } from 'rxjs';
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
                <h2>Getting actual emulation time</h2>
                @for(current_data of actual_times; track $index){
                    <li>{{current_data.id}} - {{current_data.actual_time}}</li>
                } @empty {
                    <li>There are no data</li>
                }
                <h2>Rewinding emulation</h2>
                <p>
                    <input type="datetime-local" step="1" [(ngModel)]="rewind_time">
                    <button (click)="rewindEmulation()">Set emulation time</button>
                </p>
                <h2>Manually work with registers</h2>
                <p>All fields except the data field are used for reading. For writting data ignored "count" and write data with ", "</p>
                <p>
                    <label>slave ID: </label>
                    <input type="number" min="0" [(ngModel)]="slaveID">
                    <select [(ngModel)]="registerType">
                        <option>coils</option>
                        <option>DI</option>
                        <option>HR</option>
                        <option>IR</option>
                    </select>
                    <select [(ngModel)]="operation">
                        <option value="read">read</option>
                        <option value="write">write</option>
                    </select>
                    <label>start address: </label>
                    <input type="number" min="0" [(ngModel)]="startIndex">
                    <label>count: </label>
                    <input type="number" min="1" [(ngModel)]="count">
                </p>
                <p>
                    <input [(ngModel)]="registerData">
                    <button (click)="manuallyRegistersWork()">Do action</button>
                </p>`
})

export class AppComponent{  
    actual_times: ActualTime[]=[];
    server_id: number = -1
    status: string = "Waiting"
    rewind_time: string
    constructor(private http: HttpClient) {}
    URLHead = "modbus-emulator" 
    getActualTime(): void {
        let request = `${this.URLHead}/time/actual`
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
        let request = `${this.URLHead}/time/rewind_emulation?timepoint=${this.rewind_time}`
        if (this.server_id!==-1) {
            request = request + `&server_id=${this.server_id}`
        }
        this.http.post(request, {}).subscribe({error: error => console.log(error)});
        this.status = "Success"
    }
    slaveID: number
    registerType: string
    operation: string
    startIndex: number
    count: number
    registerData: string
    manuallyRegistersWork(): void {
        let request = `${this.URLHead}/registers?server_id=${this.server_id}&slave_id=${this.slaveID}&type=${this.registerType}&start_index=${this.startIndex}`
        if (this.operation === "read") {
            request += `&count=${this.count}`
            this.http.get(request).subscribe({next:(data:any) => this.registerData=data, error: error => console.error(error)});
            return
        } else {
            this.http.post(request, {registers: JSON.parse(`[${this.registerData}]`)}).subscribe({error: error => console.log(error)});
        }
    }
    subscription = interval(500).subscribe(val => this.getActualTime());
    ngOnDestroy() {
        this.subscription.unsubscribe();
    }    
}
