BasicUpstart2(run)
*=$0810 "main"
.var ekr = $C000 //color data dest $C000-$C3E8 - Hi nibble=Fg colour; Lo nibble=Bg colour 
.var graf= $E000 //cursor pattern Hires screen 8KB the last ram blocks on C64 $E000-$FFFF
.var sc0 = 00
.var sc1 = 04
.var srcZero = ($fb)
.var destZero = ($fd)
run:
lda$d011
ora#$20
sta$d011
lda$dd00
and#$fc
sta$dd00
lda#$08
sta$d018
lda#$c0
sta$648

jsr createSurface
zap:
ldx#1
stx cordX
ldy#24
sty cordY
lda # sc0
sta buffer
jsr renderScreen
jsr setScreen0

//jsr createSurface
jsr Init // go to raster interrupts initialization


end: jmp end


    data:
    //.byte $88,$22,$88,$22,$88,$22,$88,$22
    //.byte $22,$88,$22,$88,$22,$88,$22,$88
//    .byte $00,$55,$00,$aa,$00,$55,$00,$aa
    .byte $55,$00,$aa,$00,$55,$00,$aa,$00
    //.byte $88,$44,$22,$11,$88,$44,$22,$11
    //.byte $91,$02,$41,$01,$24,$08,$12,$44
    //.byte $81,$40,$14,$29,$94,$28,$02,$90
    //.byte $55,$80,$22,$80,$01,$44,$01,$aa
    createSurface:
    lda#00
    sta$d020
    sta$d021
    
    tay
    tax
!loop:
    lda data,y
    sta lab1:graf,x
    iny
    cpy#$08
    bne !+
    ldy#00
!:
    inx
    cpx#$00
    bne !loop-
    inc lab1+1
    lda lab1+1
    cmp #$00
    bne !loop-
rts

setScreen0:
lda#$08
sta $d018
rts

setScreen1:
lda#$18
sta$d018
rts

horizontalShiftScreen0: .byte 00
horizontalShiftScreen1: .byte 00
verticalShiftScreen0: .byte 00
verticalShiftScreen1: .byte 00
animX: .byte 00
animY: .byte 00
MyIrqMainLoop:
sei
lda activeScreen: #$00
cmp # sc0        
bne showSreen0RenderScreen1 //if not equal sc0 then show sc0 and render sc1 in the background

showSreen1RenderScreen0:
    jsr setScreen1
    lda horizontalShiftScreen1
    sta $d016
    lda verticalShiftScreen1
    ora #%00110000
    sta $d011
    lda #sc1
    sta activeScreen

    //render screen 0
    inc animX
    ldx animX
    lda grubySin,x 
    sta cordX
    lda chudySin,x 
    sta horizontalShiftScreen0
    
    inc animY 
    inc animY
    inc animY
    ldy animY
    lda grubyCos,y 
    sta cordY
    lda chudyCos,y 
    sta verticalShiftScreen0
    lda #sc0
    sta buffer
jmp doTheJob

showSreen0RenderScreen1:
    //display screen 0
    jsr setScreen0
    lda horizontalShiftScreen0
    sta $d016
    lda verticalShiftScreen0
    ora #%00110000
    sta $d011
    lda #sc0
    sta activeScreen

    //render screen 1
    inc animX
    ldx animX
    lda grubySin,x 
    sta cordX
    lda chudySin,x 
    sta horizontalShiftScreen1
    
    inc animY
    inc animY
    inc animY
    ldy animY
    lda grubyCos,y 
    sta cordY
    lda chudyCos,y 
    sta verticalShiftScreen1
    lda #sc1
    sta buffer
    

doTheJob:
jsr renderScreen
asl $D019            //; acknowledge the interrupt by clearing the VIC's interrupt flag
cli
asl $d019
jmp $EA81            //; jump into KERNAL's standard interrupt service routine to handle keyboard scan, cursor display etc.


renderScreen:
    ldx#$00
!:
    ldy cordY: # 13
    lda dataSrc.hi,y 
    sta srcZero+1
    lda dataSrc.lo,y 
    //sta srcZero
    clc
    adc cordX:# 20
    sta srcZero
    lda #$00
    adc srcZero+1
    sta srcZero+1
    lda colDst.lo,x 
    sta destZero
    lda colDst.hi,x 
    clc 
    adc buffer: #sc0
    sta destZero+1
    ldy# 0
!:
    lda ($fb),y 
    sta ($fd),y 
    iny
    cpy #40 
    bne !-
    inc cordY
    inx
    cpx #25
    bne !--
    rts

Init:      
    sei                  //; set interrupt bit, make the CPU ignore interrupt requests
    lda #%01111111       //; switch off interrupt signals from CIA-1
    sta $DC0D
    and $D011            //; clear most significant bit of VIC's raster register
    sta $D011
    lda $DC0D            //; acknowledge pending interrupts from CIA-1
    lda $DD0D            //; acknowledge pending interrupts from CIA-2
    lda #$ff             //; set rasterline where interrupt shall occur
    sta $D012
    lda #<MyIrqMainLoop            //; set interrupt vectors, pointing to interrupt service routine below
    sta $0314
    lda #>MyIrqMainLoop
    sta $0315
    lda #%00000001       //; enable raster interrupt signals from VIC
    sta $D01A
    cli                  //; clear interrupt flag, allowing the CPU to respond to interrupt requests
rts





// Data imports and generators:

    .align $0100
* = * "gfx data"
grafa:  
    .import binary "demo.bin"
    //.import binary "demoG.bin"
    //.align $100



    .align $100
* = * "img data source"
dataSrc: .lohifill 50, grafa+(80*i)
 



    .align $10
* = * "colour dest screen 0"
colDst: .lohifill 25, ekr+(40*i)




    .align $10
* = * "Sine block"
grubySin: .fill 256, (160+159.5*sin(toRadians(i*360/256)))>>3
* = * "Sine shift"
chudySin: .fill 256, 7-((160+159.5*sin(toRadians(i*360/256)))& 7)




* = * "Cosine block"
grubyCos: .fill 256, (100+99.5*cos(toRadians(i*360/256)))>> 3
* = * "Cosine shift"
chudyCos: .fill 256, 7-((100+99.5*cos(toRadians(i*360/256))) & 7)