class Device {
  id: string;
  number: string;
  pushName: string;
  businessName: string;
  contacts: number;

  constructor(data: {
    id: string;
    number: string;
    push_name: string; // snake_case
    business_name: string; // snake_case
    contacts: number;
  }) {
    this.id = data.id;
    this.number = data.number;
    this.pushName = data.push_name; // Mapeia para pushName
    this.businessName = data.business_name; // Mapeia para businessName
    this.contacts = data.contacts;
  }
}

export default Device;
