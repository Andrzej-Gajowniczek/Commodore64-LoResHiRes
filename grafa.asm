//KickAssembler: 6502 source code for C64     - grafa.asm
    BasicUpstart2(run)


    *=$0810 "main"
    .var ekr = $C000 //color data dest $C000-$C3E8 - Hi nibble=Fg colour; Lo nibble=Bg colour 
    .var graf= $E000 //cursor pattern Hires screen 8KB the last ram blocks on C64 $E000-$FFFF
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
jsr czyszczenie2
jsr fillo
rts
    data:
    .byte $88,$22,$88,$22,$88,$22,$88,$22

    czyszczenie2:
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

fillo:
!:  lda lokk:grafa,x 
    sta ziel:ekr,x
    inx
    bne !-
    inc lokk+1
    inc ziel+1
    ldy ziel+1
    cpy#$C4
    bne !-
    rts
.align $0100
* = $1000 "gfx data"

grafa:  
.import binary "c64.bin"
.align $100
* = * "fill data Lo byte src"
srcLo: .fill 50, $50*i
* = * "fill data Hi byte src"
srcHi: .fill 50, grafa>>8+(0.25*i)
