package services

import (
	//addressCliente "mvc-go/clients/address"
	//telephoneCliente "mvc-go/clients/telephone"
	reservCliente "backend/clients-DAO/reserv"
	hotelCliente "backend/clients-DAO/hotel"
	"backend/model"

	//"mvc-go/clients-DAO/address"
	"backend/dto" //contienelas estructuras de datos de transferencia de objetos (DTO)
	//"backend/model"          //contiene las estructuras de datos que representan los modelos de usuario, dirección, número de teléfono,
	e "backend/utils/errors" //contiene el paquete errors
	"time"
   "math"
	"fmt"
)

type reservService struct{}
const TamHabitaciones=2; 

type reservServiceInterface interface {
	
	InsertRerserv(reservDto dto.ReservDto)(dto.ReservDto, e.ApiError)
	GetReservs()(dto.ReservsDto, e.ApiError)
	CheckAvailability(disponibilidadDto dto.CheckAvailabilityDto) (dto.CheckAvailabilityDto,e.ApiError)
	DeleteReserv(id int)(e.ApiError)
	UpdateReserv(id int, reservDto dto.ReservDto) (dto.ReservDto, e.ApiError)
	GetReservById(id int) (dto.ReservDto, e.ApiError)
	GetFutureReservsByUser(userId int) (dto.ReservsDto, e.ApiError)
	GetAllReservByHotel(id int )(dto.ReservsDto,e.ApiError)
	SearchAvaliabity( searchAvaliabity dto.RequesthHotelsAvaibylityDto)(dto.ResponseHotelsAvaibylityDtos, e.ApiError)
	
}

var ReservService reservServiceInterface

func init() {
	ReservService = &reservService{}
}


func (s *reservService) CheckAvailability(disponibilidadDto dto.CheckAvailabilityDto) (dto.CheckAvailabilityDto,e.ApiError) {

	println("pesonas antes:",disponibilidadDto.Personas )
	disponibilidadDto.Personas=s.RoomCalculation(disponibilidadDto.Personas)
    println("habitaciones que necesitan:",disponibilidadDto.Personas )

	disp, err:=reservCliente.CheckAvailability(disponibilidadDto)
   
	if err !=nil{
		fmt.Println("error aca3", e.NewBadRequestApiError(err.Error()))
		return disponibilidadDto, e.NewBadRequestApiError(err.Error())
	}

	return disp, nil

}

func (s *reservService)GetReservs()(dto.ReservsDto, e.ApiError){
	var reservsModel model.Reservs
	var ReservsDto   dto.ReservsDto

	reservsModel,err:=reservCliente.GetReservs()

	if err!=nil{
       return ReservsDto, e.NewBadRequestApiError(err.Error())
	}

	for _ ,reserv := range reservsModel{
      var reservDto dto.ReservDto
	  reservDto.Id=reserv.Id
      reservDto.UserId = reserv.UserId
	  reservDto.HotelId= reserv.HotelId
	  reservDto.DateStart=reserv.DateStart
	  reservDto.DateFinish=reserv.DateFinish
	  reservDto.DateActual=reserv.CreatedAt
      reservDto.HotelRooms=reserv.HotelRooms
	  reservDto.TotalPrice=reserv.TotalPrice

      ReservsDto = append(ReservsDto, reservDto)
	}
	return ReservsDto , nil
}

func (s *reservService)GetAllReservByHotel(id int )(dto.ReservsDto,e.ApiError){
	var reservsModel model.Reservs
	var ReservsDto   dto.ReservsDto

	reservsModel,err:=reservCliente.GetAllReservsByHotel(id)

	if err!=nil{
       return ReservsDto, e.NewBadRequestApiError(err.Error())
	}

	for _ ,reserv := range reservsModel{
      var reservDto dto.ReservDto
      reservDto.UserId = reserv.UserId
	  reservDto.HotelId= reserv.HotelId
	  reservDto.DateStart=reserv.DateStart
	  reservDto.DateFinish=reserv.DateFinish
	  reservDto.DateActual=reserv.CreatedAt
      reservDto.HotelRooms=reserv.HotelRooms
	  reservDto.TotalPrice=reserv.TotalPrice

      ReservsDto = append(ReservsDto, reservDto)
	}
	return ReservsDto , nil
}



func (s *reservService) GetFutureReservsByUser(userId int) (dto.ReservsDto, e.ApiError) {
	var reservsModel model.Reservs
	var reservsDto dto.ReservsDto

	reservsModel, err := reservCliente.GetFuturesReservsByUser(userId)
	if err != nil {
		return reservsDto, e.NewBadRequestApiError(err.Error())
	}

	for _, reserv := range reservsModel {
		var r dto.ReservDto
		r.Id = reserv.Id
		r.UserId = reserv.UserId
		r.HotelId = reserv.HotelId
		r.DateStart = reserv.DateStart
		r.DateFinish = reserv.DateFinish
		r.DateActual = reserv.CreatedAt
		r.HotelRooms = reserv.HotelRooms
		r.TotalPrice = reserv.TotalPrice
		reservsDto = append(reservsDto, r)
	}

	return reservsDto, nil
}



func (s *reservService)	InsertRerserv(reservDto dto.ReservDto)(dto.ReservDto, e.ApiError){

	//antes de hacer un insert, deberia chekear nuevamente la reserva, puede ser que entre la anterio y esta se allan acabado
    println("pesonas antes:",reservDto.HotelRooms)
	reservDto.HotelRooms=s.RoomCalculation(reservDto.HotelRooms)
    println("habitaciones que necesitan:",reservDto.HotelRooms )

	var CheckAvailability dto.CheckAvailabilityDto

	CheckAvailability.HotelId=reservDto.HotelId
	CheckAvailability.DateStart=reservDto.DateStart
	CheckAvailability.DateFinish=reservDto.DateFinish
	CheckAvailability.Personas=reservDto.HotelRooms

	disp, err:=reservCliente.CheckAvailability(CheckAvailability)

	if err!=nil {
		return reservDto, e.NewBadRequestApiError(err.Error())
	}
	if (disp.Avaliable){
       var reserv model.Reserv
	   var price int = (DiasEntre(reservDto.DateStart,reservDto.DateFinish))*reservDto.HotelRooms*120
        
	   reserv.HotelId=reservDto.HotelId
	   reserv.UserId=reservDto.UserId
	  
	   reserv.CreatedAt=reservDto.DateActual
	   reserv.DateStart=reservDto.DateStart
	   reserv.DateFinish=reservDto.DateFinish
	 
	   reserv.HotelRooms=reservDto.HotelRooms
	   reserv.TotalPrice=float32(price)
	   

	    reserv, err = reservCliente.InsertReserv(reserv)

		if err!=nil{
			return reservDto, e.NewBadRequestApiError(err.Error())
		}

		reservDto.TotalPrice=float32(price)
	
	}else{
		return reservDto, e.NewBadRequestApiError("No hay disponiblidad en la fecha seleccionada")
	}
  
     return reservDto,nil
}

func (s *reservService) DeleteReserv(id int) e.ApiError {
	err := reservCliente.DeleteReserv(id)
	if err != nil {
		return e.NewBadRequestApiError(err.Error())
	}
	return nil
}


func (s *reservService) UpdateReserv(id int, reservDto dto.ReservDto) (dto.ReservDto, e.ApiError) {
	
	//Obtengo la reserva actual
	oldReserv, err := reservCliente.FindReservById(id)
	if err != nil {
		return reservDto, e.NewNotFoundApiError("Reserva no encontrada")
	}

	//(de personas -> habitaciones)
	fmt.Println("personas (actualizadas):", reservDto.HotelRooms)
	reservDto.HotelRooms = s.RoomCalculation(reservDto.HotelRooms)
	fmt.Println("habitaciones necesarias:", reservDto.HotelRooms)

	
	var checkAvailability dto.CheckAvailabilityDto
	checkAvailability.HotelId = oldReserv.HotelId
	checkAvailability.DateStart = reservDto.DateStart
	checkAvailability.DateFinish = reservDto.DateFinish
	checkAvailability.Personas = reservDto.HotelRooms

	// Verificar disponibilidad en las nuevas fechas
	disponibilidad, errDisp := reservCliente.CheckAvailability(checkAvailability)
	if errDisp != nil {
		return reservDto, e.NewBadRequestApiError(errDisp.Error())
	}
	if !disponibilidad.Avaliable {
		return reservDto, e.NewBadRequestApiError("No hay disponibilidad para las nuevas fechas o habitaciones")
	}
	// Recalcular precio si cambio la cantidad de dias o personas
	dias := DiasEntre(reservDto.DateStart, reservDto.DateFinish)
	newPrice := float32(dias * reservDto.HotelRooms * 120)

	//Actualizar campos
	oldReserv.DateStart = reservDto.DateStart
	oldReserv.DateFinish = reservDto.DateFinish
	oldReserv.HotelRooms = reservDto.HotelRooms
	oldReserv.TotalPrice = newPrice

	errUpdate := reservCliente.UpdateReserv(oldReserv)
	if errUpdate != nil {
		return reservDto, e.NewInternalServerApiError("Error al actualizar la reserva", errUpdate)
	}


	reservDto.TotalPrice = newPrice
	reservDto.UserId = oldReserv.UserId
	reservDto.HotelId = oldReserv.HotelId
	reservDto.DateActual = oldReserv.CreatedAt 

	return reservDto, nil
}

func (s *reservService) GetReservById(id int) (dto.ReservDto, e.ApiError) {
	var reservDto dto.ReservDto

	reservModel, err := reservCliente.FindReservById(id)
	if err != nil {
		return reservDto, e.NewNotFoundApiError("Reserva no encontrada")
	}
    reservDto.Id=reservModel.Id
	reservDto.UserId = reservModel.UserId
	reservDto.HotelId = reservModel.HotelId
	reservDto.DateStart = reservModel.DateStart
	reservDto.DateFinish = reservModel.DateFinish
	reservDto.DateActual = reservModel.CreatedAt
	reservDto.HotelRooms = reservModel.HotelRooms
	reservDto.TotalPrice = reservModel.TotalPrice

	return reservDto, nil
}


func (s *reservService) SearchAvaliabity(searchAvaliabity dto.RequesthHotelsAvaibylityDto)(dto.ResponseHotelsAvaibylityDtos, e.ApiError){

	var responseHotels dto.ResponseHotelsAvaibylityDtos


   //(de personas -> habitaciones)
   searchAvaliabity.Personas = s.RoomCalculation(searchAvaliabity.Personas)
    
   //pregunto por todos los ids de los hoteles que estan disponibles
   ids, err := reservCliente.GetHotelesDisponibles(searchAvaliabity)
  if(err!=nil){
	return   responseHotels, e.NewNotFoundApiError("error en la busqueda de hoteles")
  }

   //consulto los datos de los hoteles correspondientes a dichos IDs:

   hotels,err := hotelCliente.GetHotelsById(ids)
  if(err!=nil){
	return responseHotels,e.NewNotFoundApiError("error en la busqueda de hoteles")
  }
   
   for _, hotel := range hotels {
		var resposeHotel dto.ResponseHotelAvaibylityDto
		resposeHotel.Id = hotel.Id
		resposeHotel.Name = hotel.Name
		resposeHotel.Description=hotel.Description
		resposeHotel.Rooms = hotel.Rooms

		var amenityDtos []dto.AmenitiDto
		for _, amenity := range hotel.Amenities {
			amenityDto := dto.AmenitiDto{
				Id: amenity.Id,
				Name:        amenity.Name,
				Description: amenity.Description,
			}
			amenityDtos = append(amenityDtos, amenityDto)
		}
		var imagenesDto []dto.ImagenDto
		for _, imagen := range hotel.Imagenes {
			imagenDto := dto.ImagenDto{
				Url: imagen.Url,
				Id:  imagen.Id,
			}
			imagenesDto = append(imagenesDto, imagenDto)
		}

		resposeHotel.Imagenes = imagenesDto
		resposeHotel.Amenities = amenityDtos

		responseHotels = append(responseHotels, resposeHotel)
	}

	return responseHotels,nil
}




//funciones extra: 


func DiasEntre(fechaInicio, fechaFin time.Time) int {
	inicio := time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, fechaInicio.Location())
	fin := time.Date(fechaFin.Year(), fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, fechaFin.Location())

	// Calcula la diferencia y convierte a días
	duracion := fin.Sub(inicio)
	dias := int(math.Ceil(duracion.Hours() / 24))

	// (si la fecha de inicio está después de la de fin)
	if dias < 0 {
		dias = -dias
	}

	return dias
}


func (s *reservService) RoomCalculation (personas int)(habitaciones int){
		if((personas%TamHabitaciones)==0){
	   habitaciones=personas/TamHabitaciones
	}else{
       habitaciones= (int(personas/2))+1
	}
	return habitaciones
}

